package models

import (
	"strings"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	//Password represents a user password in database
	Password struct {
		ID        string        `bson:"_id" json:"password_id"`
		Password  []byte        `bson:"password" json:"-"`
		Hasher    string        `bson:"hasher" json:"-"`
		Confirmed bool          `bson:"confirmed" json:"-"`
		UserID    string        `bson:"user_id" json:"user_id" validate:"required"`
		Modified  vmod.Modified `bson:"modified" json:"-"`
	}
	//PasswordReset represents a request for providing password reset.
	PasswordReset struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

//NewPassword initial Password struct
func NewPassword(p string, userID string, cTime int64) (pw *Password, apiErr *verr.APIError) {
	password, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	pw = &Password{
		ID:       uuid.New().String(),
		UserID:   userID,
		Password: password,
		Hasher:   "bcrypt",
		Modified: vmod.Modified{
			Created: cTime,
			Updated: cTime,
		},
		Confirmed: false,
	}
	return pw, nil
}

//Validate provides a bcrypt validation for p.
func (p *Password) Validate(v string) (apiErr *verr.APIError) {
	err := bcrypt.CompareHashAndPassword(p.Password, []byte(v))
	if err != nil {
		if strings.Contains(err.Error(), "is not the hash of the given password") {
			return verr.NewAPIError(err).Forbidden("wrong_password")
		}
		return verr.NewAPIError(err).InternalServerError()
	}
	return nil
}
