package dao

import (
	"context"
	"drops-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/verr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateApplication creates application model in database
func CreateApplication(a models.ApplicationCreate) (*models.Application, *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app := a.Application()
	id, err := DB.Collection("applications").InsertOne(ctx, app)
	if err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	app.ID = id.InsertedID.(primitive.ObjectID)
	return app, nil
}

//GetApplicationList return all applications
func GetApplicationList() (aList []models.Application, apiErr *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := DB.Collection("applications").Find(ctx, bson.M{})
	if err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	if err = cursor.All(ctx, &aList); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	return aList, nil
}
