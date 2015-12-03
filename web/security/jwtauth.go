package security

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alinz/releasifier/lib/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

var (
	errUnauthorized = errors.New("unauthorized token")
)

type JwtAuth struct {
	signKey   []byte
	verifyKey []byte
	signer    jwt.SigningMethod
}

// verifyKey is only for RSA
func New(alg string, signKey []byte, verifyKey []byte) *JwtAuth {
	return &JwtAuth{
		signKey:   signKey,
		verifyKey: verifyKey,
		signer:    jwt.GetSigningMethod(alg),
	}
}

func (ja *JwtAuth) Handle(paramAliases ...string) func(chi.Handler) chi.Handler {
	return func(next chi.Handler) chi.Handler {
		hfn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

			var tokenStr string
			var err error

			// Get token from query params
			tokenStr = r.URL.Query().Get("jwt")

			// Get token from other query param aliases
			if tokenStr == "" && paramAliases != nil && len(paramAliases) > 0 {
				for _, p := range paramAliases {
					tokenStr = r.URL.Query().Get(p)
					if tokenStr != "" {
						break
					}
				}
			}

			// Get token from authorization header
			if tokenStr == "" {
				bearer := r.Header.Get("Authorization")
				if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
					tokenStr = bearer[7:]
				}
			}

			// Get token from cookie
			if tokenStr == "" {
				cookie, err := r.Cookie("jwt")
				if err == nil {
					tokenStr = cookie.Value
				}
			}

			// Token is required, cya
			if tokenStr == "" {
				err = errUnauthorized
			}

			// Verify the token
			token, err := ja.Decode(tokenStr)
			if err != nil || !token.Valid || token.Method != ja.signer {
				utils.Respond(w, 401, errUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, "jwt", token.Raw)
			ctx = context.WithValue(ctx, "jwt.token", token)

			next.ServeHTTPC(ctx, w, r)
		}
		return chi.HandlerFunc(hfn)
	}
}

func (ja *JwtAuth) Handler(next chi.Handler) chi.Handler {
	return ja.Handle("")(next)
}

func (ja *JwtAuth) Encode(claims map[string]interface{}) (t *jwt.Token, tokenString string, err error) {
	t = jwt.New(ja.signer)
	t.Claims = claims
	tokenString, err = t.SignedString(ja.signKey)
	t.Raw = tokenString
	return
}

func (ja *JwtAuth) Decode(tokenString string) (t *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if ja.verifyKey != nil && len(ja.verifyKey) > 0 {
			return ja.verifyKey, nil
		} else {
			return ja.signKey, nil
		}
	})
}
