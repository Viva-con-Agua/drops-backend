package dao

import (
	"context"
	"drops-backend/models"
	"time"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

func SignUp(s *models.SignUp) (user *vmod.User, api_err *verr.ApiError) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c_time := time.Now().Unix()

	user = s.SignUser.NewUser(c_time)
	if api_err := UserInsertOne(ctx, user); api_err != nil {
		return nil, api_err
	}

	profile := s.SignUser.NewProfile(c_time, user.ID)
	if api_err := ProfileInsertOne(ctx, profile); api_err != nil {
		return nil, api_err
	}

	if password, api_err := models.InitPassword(s.SignUser.Password, user.ID, c_time); api_err != nil {
		return nil, api_err
	} else if api_err := PasswordInsertOne(ctx, password); api_err != nil {
		return nil, api_err
	}

	if token, api_err := vmod.InitToken("sign_up", c_time, time.Hour*24*7, user.ID); api_err != nil {
		return nil, api_err
	} else if api_err := TokenInsertOne(ctx, token); api_err != nil {
		return nil, api_err
	}

	user.Profile = *profile

	return user, verr.GetApiError(nil, nil)
}

func SignIn(s *models.SignIn) (user *vmod.User, api_err *verr.ApiError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": s.Email}
	user, api_err = UserFindOne(ctx, filter)
	if api_err != nil {
		return nil, api_err
	}
	password := new(models.Password)
	filter = bson.M{"user_id": user.ID}
	password, api_err = PasswordFindOne(ctx, filter)
	if api_err != nil {
		return nil, api_err
	}

	api_err = password.Validate(s.Password)
	if api_err != nil {
		return nil, api_err
	}

	filter = bson.M{"user_id": user.ID}
	profile, api_err := ProfileFindOne(ctx, filter)
	if api_err != nil {
		return nil, api_err
	}
	user.Profile = *profile
	return user, nil
}
