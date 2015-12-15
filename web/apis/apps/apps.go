package apps

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pressly/selfie/data"
	"github.com/pressly/selfie/errors"
	"github.com/pressly/selfie/lib/utils"
	"github.com/pressly/selfie/web/constants"
	"github.com/pressly/selfie/web/security"
	"github.com/pressly/selfie/web/util"
	"golang.org/x/net/context"
)

func getAllApps(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)

	apps, err := data.DB.App.FindAllApps(userID)

	if err != nil {
		utils.Respond(w, 400, err)
	} else {
		utils.Respond(w, 200, apps)
	}
}

func getApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")

	app, err := data.DB.App.FindApp(appID, userID)

	if err == nil {
		utils.Respond(w, 200, app)
	} else {
		utils.Respond(w, 400, err)
	}
}

func generateAppToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")
	generateAppTokenReq := ctx.Value(constants.CtxKeyParsedBody).(*generateAppTokenRequest)

	if !data.DB.App.HasPermission(appID, userID, data.ADMIN, data.OWNER) {
		utils.RespondEx(w, nil, 0, errors.ErrorAuthorizeAccess)
		return
	}

	permission, err := data.GetPermssionString(*generateAppTokenReq.Permission)

	if err != nil {
		utils.RespondEx(w, nil, 0, errors.ErrorSomethingWentWrong)
		return
	}

	claims := map[string]interface{}{
		"app_id":     fmt.Sprintf("%v", appID),
		"permission": fmt.Sprintf("%v", permission),
	}
	_, tokenStr, err := security.TokenAuth.Encode(claims)

	if err != nil {
		utils.RespondEx(w, nil, 0, errors.ErrorSomethingWentWrong)
		return
	}

	utils.RespondEx(w, generateAppTokenResponse{Token: tokenStr}, http.StatusOK, nil)
}

func acceptAppToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")
	appTokenReq := ctx.Value(constants.CtxKeyParsedBody).(*appTokenRequest)

	//decode jwt token
	token, err := security.TokenAuth.Decode(*appTokenReq.Token)
	if err != nil || !token.Valid {
		utils.RespondEx(w, nil, 401, errors.ErrorAuthorizeAccess)
		return
	}

	tokenAppID, err := strconv.ParseInt(token.Claims["app_id"].(string), 10, 64)
	if err != nil || tokenAppID != appID {
		utils.RespondEx(w, nil, 402, errors.ErrorAuthorizeAccess)
		return
	}

	tokenPermission, err := data.GetPermissionByName(token.Claims["permission"].(string))
	if err != nil || tokenPermission == data.ANONYMOUSE {
		utils.RespondEx(w, nil, 403, errors.ErrorAuthorizeAccess)
		return
	}

	//check if user has already have an access
	if data.DB.App.HasPermission(appID, userID, data.ADMIN, data.OWNER, data.MEMBER) {
		utils.RespondEx(w, nil, 404, errors.ErrorAlreadyAcceessed)
		return
	}

	//try to grand access to app with authorized permission
	if !data.DB.App.GrantAccess(appID, userID, tokenPermission) {
		utils.RespondEx(w, nil, 405, errors.ErrorAppNotFound)
		return
	}

	utils.Respond(w, 200, nil)
}

func createApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	createAppReq := ctx.Value(constants.CtxKeyParsedBody).(*createAppRequest)

	app, err := data.DB.App.CreateNewApp(1, *createAppReq.Name)

	if err == nil {
		utils.Respond(w, 200, app)
	} else {
		utils.Respond(w, 400, err)
	}
}

func updateApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")
	updateAppReq := ctx.Value(constants.CtxKeyParsedBody).(*updateAppRequest)

	err := data.DB.App.UpdateApp(appID, updateAppReq.Name, updateAppReq.PublicKey, updateAppReq.PrivateKey, updateAppReq.Private, userID)

	if err == nil {
		utils.Respond(w, 200, nil)
	} else {
		utils.Respond(w, 400, err)
	}
}

func removeApp(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")

	err := data.DB.App.RemoveApp(appID, userID)
	if err == nil {
		utils.Respond(w, 200, nil)
	} else {
		utils.Respond(w, 400, err)
	}
}
