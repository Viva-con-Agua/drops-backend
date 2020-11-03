package dao

import (
	"context"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
)

func TokenInsertOne(ctx context.Context, t *vmod.Token) (api_err *verr.ApiError) {
	token_col := utils.Database.Collection("token")
	if _, err := token_col.InsertOne(ctx, t); err != nil {
		return verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	return nil
}
