package dao

import (
	"context"
	"drops-backend/models"
	"drops-backend/utils"
	"net/http"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"go.mongodb.org/mongo-driver/bson"
)

func PasswordInsertOne(ctx context.Context, p *models.Password) (api_err *verr.ApiError) {

	var password_col = utils.Database.Collection("password")
	if _, err := password_col.InsertOne(ctx, p); err != nil {
		return verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return nil
}

func PasswordFindOne(ctx context.Context, filter bson.M) (password *models.Password, api_err *verr.ApiError) {

	var password_col = utils.Database.Collection("password")
	password = new(models.Password)
	err := password_col.FindOne(ctx, filter).Decode(&password)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			r_err := verr.ResponseError{
				Code: http.StatusNotFound,
				Response: verr.ResponseMessage{
					Message: "password_not_found",
				},
			}
			return nil, verr.GetApiError(err, &r_err)
		}
		return nil, verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return password, nil
}
