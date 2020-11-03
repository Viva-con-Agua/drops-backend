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

func UserInsertOne(ctx context.Context, u *vmod.User) (api_err *verr.ApiError) {

	var user_col = utils.Database.Collection("user")
	if _, err := user_col.InsertOne(ctx, u); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			r_err := verr.ResponseError{
				Code: http.StatusConflict,
				Response: verr.ResponseMessage{
					Message: "email_duplicate",
				},
			}
			return verr.GetApiError(err, &r_err)
		}
		return verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return nil
}

func UserFindOne(ctx context.Context, filter bson.M) (user *vmod.User, api_err *verr.ApiError) {

	var user_col = utils.Database.Collection("user")
	user = new(vmod.User)
	err := user_col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			r_err := verr.ResponseError{
				Code: http.StatusNotFound,
				Response: verr.ResponseMessage{
					Message: "user_not_found",
				},
			}
			return nil, verr.GetApiError(err, &r_err)
		}
		return nil, verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return user, nil
}
