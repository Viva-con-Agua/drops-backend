package dao

import (
	"context"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

//UserInsertOne inserts one User model into database
func UserInsertOne(ctx context.Context, u *vmod.User) (apiErr *verr.APIError) {
	var coll = DB.Collection("user")
	if _, err := coll.InsertOne(ctx, u); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return verr.NewAPIError(err).Conflict("email_duplicate")
		}
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}

//UserFindOne selects a User model from database by using `filter`
func UserFindOne(ctx context.Context, filter bson.M) (user *vmod.User, apiErr *verr.APIError) {
	var coll = DB.Collection("user")
	user = new(vmod.User)
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("user_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	return user, nil
}
