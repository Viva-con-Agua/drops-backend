package dao

import (
	"context"
	"drops-backend/utils"
	"net/http"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

func ProfileInsertOne(ctx context.Context, p *vmod.Profile) (api_err *verr.ApiError) {

	var profile_col = utils.Database.Collection("profile")
	if _, err := profile_col.InsertOne(ctx, p); err != nil {
		return verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return nil
}

func ProfileFindOne(ctx context.Context, filter bson.M) (profile *vmod.Profile, api_err *verr.ApiError) {

	var profile_col = utils.Database.Collection("profile")
	profile = new(vmod.Profile)
	err := profile_col.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			r_err := verr.ResponseError{
				Code: http.StatusNotFound,
				Response: verr.ResponseMessage{
					Message: "profile_not_found",
				},
			}
			return nil, verr.GetApiError(err, &r_err)
		}
		return nil, verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return profile, nil
}
