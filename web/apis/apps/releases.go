package apps

import (
	"net/http"

	"github.com/pressly/selfie/data"
	"github.com/pressly/selfie/lib/utils"
	"github.com/pressly/selfie/web/constants"
	"github.com/pressly/selfie/web/util"
	"golang.org/x/net/context"
)

func getAllReleases(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//get userID and appID
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")

	releases, err := data.DB.Release.FindAllReleases(userID, appID)

	utils.RespondEx(w, releases, 0, err)
}

func updateRelease(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")
	releaseID, _ := util.GetParamValueAsID(ctx, "releaseID")

	//grabing update release request
	updateReleaseReq := ctx.Value(constants.CtxKeyParsedBody).(*updateReleaseRequest)

	err := data.DB.Release.UpdateRelease(updateReleaseReq.Note, updateReleaseReq.Platform, updateReleaseReq.Version, releaseID, appID, userID)
	utils.RespondEx(w, nil, 0, err)
}

func createRelease(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//get userID and appID
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")

	//grabing release request
	createReleaseReq := ctx.Value(constants.CtxKeyParsedBody).(*createReleaseRequest)

	//try to create release and return created release record
	release, err := data.DB.Release.CreateRelease(*createReleaseReq.Version, *createReleaseReq.Platform, createReleaseReq.Note, userID, appID)
	if err == nil {
		utils.Respond(w, 200, release)
	} else {
		utils.Respond(w, 400, err)
	}
}

func lockRelease(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID, _ := util.GetUserIDFromContext(ctx)
	appID, _ := util.GetParamValueAsID(ctx, "appID")
	releaseID, _ := util.GetParamValueAsID(ctx, "releaseID")

	err := data.DB.Release.LockRelease(releaseID, appID, userID)
	utils.RespondEx(w, nil, 0, err)
}
