package dao

import (
	"context"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
)

//TokenInsertOne inserts token model into database
func TokenInsertOne(ctx context.Context, t *vmod.Token) (api_err *verr.APIError) {
	coll := DB.Collection("token")
	if _, err := coll.InsertOne(ctx, t); err != nil {
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}
