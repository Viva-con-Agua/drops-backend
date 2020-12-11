package dao

import (
	"context"
	"drops-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

//GetProfileList return a list of Profiles
func GetProfileList(q models.ProfileQuery) (pList []vmod.Profile, apiErr *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := DB.Collection("profile").Find(ctx, q.Filter())
	if err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	if err = cursor.All(ctx, &pList); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	return pList, nil
}

//UpdateProfile updates the profile model in database
func UpdateProfile(p *vmod.Profile) (*vmod.Profile, *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := DB.Collection("profile").UpdateOne(ctx, bson.M{"user_id": p.UserID}, bson.M{"$set": p}); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	return p, nil
}
