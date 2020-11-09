package dao

import (
	"context"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

func ProfileInsertOne(ctx context.Context, p *vmod.Profile) (api_err *verr.APIError) {
	var coll = DB.Collection("profile")
	if _, err := coll.InsertOne(ctx, p); err != nil {
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}

func ProfileFindOne(ctx context.Context, filter bson.M) (profile *vmod.Profile, api_err *verr.APIError) {
	var coll = DB.Collection("profile")
	profile = new(vmod.Profile)
	err := coll.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("profile_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	return profile, nil
}
