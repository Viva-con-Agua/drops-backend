package models

import (
	"time"

	"github.com/Viva-con-Agua/echo-pool/crm"
	"github.com/Viva-con-Agua/vcago/civi"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
)

type (
	// for /auth/signin
	SignIn struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Service  string `json:"service" validate:"required"`
	}
	SignUser struct {
		Email         string `json:"email" validate:"required"`
		Password      string `json:"password" validate:"required"`
		FirstName     string `json:"first_name" validate:"required"`
		LastName      string `json:"last_name" validate:"required"`
		PrivacyPolicy bool   `json:"privacy_policy"`
		Country       string `json:"country"`
		Service       string `json:"service"`
	}
	Campaign struct {
		CampaignId   int    `json:"campaign_id" validate:"required"`
		CampaignName string `json:"campaign_name" validate:"required"`
	}
	SignUp struct {
		SignUser    SignUser    `json:"sign_user" validate:"required"`
		Campaign    Campaign    `json:"campaign" validate:"required"`
		Offset      civi.Offset `json:"offset" validate:"required"`
		RedirectUrl string      `json:"redirect_url" validate:"required"`
	}
	NewToken struct {
		Campaign Campaign `json:"campaign" validate:"required"`
		Email    string   `json:"email" validate:"required"`
	}
)

func (signup_data *SignUp) CrmUserSignUp(drops_id string, link string) *civi.CrmUserSignUp {
	crm_user := new(civi.CrmUserSignUp)
	crm_user.CrmData.DropsId = drops_id
	crm_user.CrmData.CampaignId = signup_data.Campaign.CampaignId
	crm_user.CrmData.Activity = "DROPS_USER_CREATE"
	crm_user.CrmData.Created = time.Now().Unix()
	crm_user.CrmUser.Email = signup_data.SignUser.Email
	crm_user.CrmUser.FirstName = signup_data.SignUser.FirstName
	crm_user.CrmUser.LastName = signup_data.SignUser.LastName
	crm_user.CrmUser.PrivacyPolicy = signup_data.SignUser.PrivacyPolicy
	crm_user.CrmUser.Country = signup_data.SignUser.Country
	crm_user.Mail.Link = link
	crm_user.Offset = signup_data.Offset
	return crm_user

}

func (nt *NewToken) CrmEmailBody(a string) *crm.CrmEmailBody {
	ce := new(crm.CrmEmailBody)
	ce.CrmData.CampaignId = nt.Campaign.CampaignId
	ce.CrmData.Activity = a
	ce.Mail.Email = nt.Email
	return ce
}

func (u *SignUser) NewUser(c_time int64) (n_u *vmod.User) {
	//create password models based on bcrypt
	modified := vmod.InitModified(c_time)
	// create bcrypt.User
	user := new(vmod.User)
	user.ID = uuid.New().String()
	user.Email = u.Email
	user.Policies = *vmod.InitPolicies("confirmed", false, c_time)
	user.Policies.Add("privacy_policy", u.PrivacyPolicy, c_time)
	user.Permission = *vmod.InitPermission("joined", "drops-service", modified.Created)
	if u.Service != "" && u.Service != "drops-service" {
		user.Permission.Add("member", u.Service, modified.Created)
	}
	user.Modified = *modified
	return user
}

func (u *SignUser) NewProfile(c_time int64, user_id string) *vmod.Profile {
	// create Profile
	modified := vmod.InitModified(c_time)
	profile := new(vmod.Profile)
	profile.ID = uuid.New().String()
	profile.UserID = user_id
	profile.FirstName = u.FirstName
	profile.LastName = u.LastName
	profile.FullName = u.FirstName + " " + u.LastName
	profile.Country = u.Country
	profile.Modified = *modified
	return profile
}
