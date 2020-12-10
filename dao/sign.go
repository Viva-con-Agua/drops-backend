package dao

import (
	"context"
	"drops-backend/models"
	"errors"
	"strings"
	"time"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson"
)

//SignUp is a mongodb controller that manages the SignUp database process.
func SignUp(s *models.SignUp) (r *models.SignUpReturn, apiErr *verr.APIError) {
	r = new(models.SignUpReturn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cTime := time.Now().Unix()

	//insert user model
	user := s.SignUser.User(cTime)
	if _, err := DB.Collection("user").InsertOne(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, verr.NewAPIError(err).Conflict("email_duplicate")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//insert profile
	profile := s.SignUser.Profile(cTime, user.ID)
	if _, err := DB.Collection("profile").InsertOne(ctx, profile); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//insert password
	if password, apiErr := models.NewPassword(s.SignUser.Password, user.ID, cTime); apiErr != nil {
		return nil, apiErr
	} else if _, err := DB.Collection("password").InsertOne(ctx, password); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//insert policies
	policies := s.SignUser.Policies(cTime, user.ID)
	if _, err := DB.Collection("policies").InsertOne(ctx, policies); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	user.Policies = *policies

	if token, apiErr := vmod.NewToken("sign_up", cTime, time.Hour*24*7, user.ID); apiErr != nil {
		return nil, apiErr
	} else if _, err := DB.Collection("token").InsertOne(ctx, token); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	} else {
		r.Token = *token
	}

	user.Profile = *profile
	r.User = *user
	return r, nil
}

//SignIn is a mongodb controller that manages the SignIn database process.
func SignIn(s *models.SignIn) (user *vmod.User, apiErr *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//find user
	filter := bson.M{"email": s.Email}
	user = new(vmod.User)
	err := DB.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("user_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//get password and validate password
	password := new(models.Password)
	filter = bson.M{"user_id": user.ID}
	password = new(models.Password)
	err = DB.Collection("password").FindOne(ctx, filter).Decode(&password)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("password_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()

	}
	apiErr = password.Validate(s.Password)
	if apiErr != nil {
		return nil, verr.NewAPIError(err).Forbidden("password_not_valid")
	}

	//get profile
	profile := new(vmod.Profile)
	err = DB.Collection("profile").FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("profile_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//get and validate policies
	policies := new(vmod.Policies)
	err = DB.Collection("policies").FindOne(ctx, filter).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("policies_not_found")
		}

		return nil, verr.NewAPIError(err).InternalServerError()
	}

	apiErr = policies.Validate("confirmed")
	if apiErr != nil {

		return nil, verr.NewAPIError(errors.New("not_confirmed")).Forbidden("not_confirmed")
	}

	user.Policies = *policies
	user.Profile = *profile
	return user, nil
}

//SignUpConfirm confirmes the user and return vmod.User for signin after confirm
func SignUpConfirm(code string) (user *vmod.User, apiErr *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cTime := time.Now().Unix()

	// get token
	filter := bson.M{"code": code, "tcase": "sign_up"}
	token := new(vmod.Token)
	err := DB.Collection("token").FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("token_not_valid")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	filter = bson.M{"user_id": token.ModelID}

	//get and confirm policies
	policies := new(vmod.Policies)
	err = DB.Collection("policies").FindOne(ctx, filter).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("policies_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	policies = policies.Update("confirmed", cTime, true)
	_, err = DB.Collection("policies").UpdateOne(ctx, bson.M{"_id": policies.ID}, bson.M{"$set": policies})
	if err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//get user model
	user = new(vmod.User)
	err = DB.Collection("user").FindOne(ctx, bson.M{"_id": token.ModelID}).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("user_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//insert profile
	profile := new(vmod.Profile)
	if _, err := DB.Collection("profile").InsertOne(ctx, profile); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	if _, err := DB.Collection("token").DeleteOne(ctx, bson.M{"_id": token.ID}); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	user.Policies = *policies
	user.Profile = *profile
	return user, nil

}

//NewSignUpToken creates new token for the sign up confirm process
func NewSignUpToken(email string) (*vmod.Token, *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cTime := time.Now().Unix()
	//find user
	filter := bson.M{"email": email}
	user := new(vmod.User)
	err := DB.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("user_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}

	//get and confirm policies
	policies := new(vmod.Policies)
	err = DB.Collection("policies").FindOne(ctx, bson.M{"user_id": user.ID}).Decode(&policies)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, verr.NewAPIError(err).NotFound("policies_not_found")
		}
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	apiErrT := policies.Validate("confirmed")
	if apiErrT == nil {
		return nil, verr.NewAPIError(errors.New("user_already_confirmed")).Forbidden("user_already_confirmed")
	}
	if _, err := DB.Collection("token").DeleteOne(ctx, bson.M{"model_id": user.ID, "t_case": "sign_up"}); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	}
	if token, apiErr := vmod.NewToken("sign_up", cTime, time.Hour*24*7, user.ID); apiErr != nil {
		return nil, apiErr
	} else if _, err := DB.Collection("token").InsertOne(ctx, token); err != nil {
		return nil, verr.NewAPIError(err).InternalServerError()
	} else {
		return token, nil
	}

}
