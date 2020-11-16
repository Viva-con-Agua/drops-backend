package dao

import (
	"context"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

//TokenInsertOne inserts token model into database
func TokenInsertOne(ctx context.Context, t *vmod.Token) (api_err *verr.APIError) {
	coll := DB.Collection("token")
	if _, err := coll.InsertOne(ctx, t); err != nil {
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}

//TokenFindOne select one policies model from database
func TokenFindOne(ctx context.Context, filter bson.M) (policies *vmod.Policies, apiErr *verr.APIError) {
	var coll = DB.Collection("policies")
	policies = new(vmod.Policies)
	err := coll.FindOne(ctx, filter).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("token_not_found")
		}
		return nil, apiErr.InternalServerError()
	}
	return policies, nil

}
