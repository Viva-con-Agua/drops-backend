package dao

import (
	"context"
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

//PoliciesCreateInit creates policies for user create or signup
//NOT USED
func PoliciesCreateInit(ctx context.Context, userID string, cTime int64) (policies *vmod.Policies, apiErr *verr.APIError) {
	var coll = DB.Collection("policies")
	policies = new(vmod.Policies)
	err := coll.FindOne(ctx, bson.M{"_id": "default_policies"}).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("default_policies_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()

	}
	policies.ID = uuid.New().String()
	policies.UserID = userID
	return policies, nil
}

//PoliciesInsertOne inserts one Policies model into database.
func PoliciesInsertOne(ctx context.Context, i *vmod.Policies) (apiErr *verr.APIError) {
	var coll = DB.Collection("policies")
	if _, err := coll.InsertOne(ctx, i); err != nil {
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}

//PoliciesFindOne select one policies model from database
func PoliciesFindOne(ctx context.Context, filter bson.M) (policies *vmod.Policies, apiErr *verr.APIError) {
	var coll = DB.Collection("policies")
	policies = new(vmod.Policies)
	err := coll.FindOne(ctx, filter).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("policies_not_found")
		}
		return nil, apiErr.InternalServerError()
	}
	return policies, nil

}
