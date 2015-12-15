package web

import (
	"net/http"

	"github.com/pressly/selfie/web/apis/apps"
	"github.com/pressly/selfie/web/apis/auth"
	"github.com/pressly/chi"
)

//New create all routers under one Handler
func New() http.Handler {
	r := chi.NewRouter()

	r.Mount("/auth", auth.Routes())
	r.Mount("/apps", apps.Routes())

	return r
}
