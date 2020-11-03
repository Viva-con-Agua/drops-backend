package models

import (
	"net/http"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	Password struct {
		ID        string        `bson:"_id" json:"password_id"`
		Password  []byte        `bson:"password" json:"-"`
		Hasher    string        `bson:"hasher" json:"-"`
		Confirmed bool          `bson:"confirmed" json:"-"`
		UserID    string        `bson:"user_id" json:"user_id" validate:"required"`
		Modified  vmod.Modified `bson:"modified" json:"-"`
	}
	PasswordReset struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

func InitPassword(p string, user_id string, c_time int64) (pw *Password, api_err *verr.ApiError) {
	password, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return nil, verr.GetApiError(err, &verr.RespErrorInternalServer)
	}
	pw = &Password{
		ID:       uuid.New().String(),
		UserID:   user_id,
		Password: password,
		Hasher:   "bcrypt",
		Modified: vmod.Modified{
			Created: c_time,
			Updated: c_time,
		},
		Confirmed: false,
	}
	return pw, api_err
}

func (p *Password) Validate(p_string string) (api_err *verr.ApiError) {
	err := bcrypt.CompareHashAndPassword(p.Password, []byte(p_string))
	if err != nil {
		r_err := verr.ResponseError{
			Code: http.StatusUnauthorized,
			Response: verr.ResponseMessage{
				Message: "password_not_valid",
			},
		}
		return verr.GetApiError(err, &r_err)
	}
	return nil
}
