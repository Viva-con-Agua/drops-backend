package dao

import (
	"context"
	"drops-backend/models"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"go.mongodb.org/mongo-driver/bson"
)

//PasswordInsertOne inserts a Password models into database.
func PasswordInsertOne(ctx context.Context, i *models.Password) (apiErr *verr.APIError) {
	var coll = DB.Collection("password")
	if _, err := coll.InsertOne(ctx, i); err != nil {
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}

//PasswordFindOne selects a Password model from database by using `filter`
func PasswordFindOne(ctx context.Context, filter bson.M) (password *models.Password, apiErr *verr.APIError) {
	var coll = DB.Collection("password")
	password = new(models.Password)
	err := coll.FindOne(ctx, filter).Decode(&password)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("password_not_found")
		}
		return nil, apiErr.InternalServerError()
	}
	return password, nil
}
