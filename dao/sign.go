package dao

import (
	"context"
	"drops-backend/models"
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

	user := s.SignUser.User(cTime)
	if apiErr := UserInsertOne(ctx, user); apiErr != nil {
		return nil, apiErr
	}

	profile := s.SignUser.Profile(cTime, user.ID)
	if apiErr := ProfileInsertOne(ctx, profile); apiErr != nil {
		return nil, apiErr
	}

	if password, apiErr := models.NewPassword(s.SignUser.Password, user.ID, cTime); apiErr != nil {
		return nil, apiErr
	} else if apiErr := PasswordInsertOne(ctx, password); apiErr != nil {
		return nil, apiErr
	}

	policies := s.SignUser.Policies(cTime, user.ID)
	if apiErr := PoliciesInsertOne(ctx, policies); apiErr != nil {
		return nil, apiErr
	}

	user.Policies = *policies
	if token, apiErr := vmod.NewToken("sign_up", cTime, time.Hour*24*7, user.ID); apiErr != nil {
		return nil, apiErr
	} else if apiErr := TokenInsertOne(ctx, token); apiErr != nil {
		return nil, apiErr
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

	filter := bson.M{"email": s.Email}
	user, apiErr = UserFindOne(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}
	password := new(models.Password)
	filter = bson.M{"user_id": user.ID}
	password, apiErr = PasswordFindOne(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	apiErr = password.Validate(s.Password)
	if apiErr != nil {
		return nil, apiErr
	}

	profile, apiErr := ProfileFindOne(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	policies, apiErr := PoliciesFindOne(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}
	apiErr = policies.Validate("confirmed")
	if apiErr != nil {
		return nil, apiErr
	}
	user.Policies = *policies
	user.Profile = *profile
	return user, nil
}

/*func SignUpConfirm(t string) (u *vmod.User, apiErr *verr.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"code": t}
	token, apiErr := TokenFindOne(ctx, filter)

}*/
