package models

import (
	"time"

	"github.com/Viva-con-Agua/echo-pool/crm"
	"github.com/Viva-con-Agua/vcago/civi"
	"github.com/Viva-con-Agua/vcago/vmod"
	"github.com/google/uuid"
)

type (
	//SignIn represents the user sign in data
	SignIn struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Service  string `json:"service" validate:"required"`
	}
	//SignUser represents the initial user data we need for sign up
	SignUser struct {
		Email         string `json:"email" validate:"required"`
		Password      string `json:"password" validate:"required"`
		FirstName     string `json:"first_name" validate:"required"`
		LastName      string `json:"last_name" validate:"required"`
		PrivacyPolicy bool   `json:"privacy_policy"`
		Country       string `json:"country"`
		Service       string `json:"service"`
	}
	//Campaign represents civi crm campaign
	Campaign struct {
		CampaignID   int    `json:"campaign_id" validate:"required"`
		CampaignName string `json:"campaign_name" validate:"required"`
	}
	//SignUp represents data collection for sign up
	SignUp struct {
		SignUser    SignUser    `json:"sign_user" validate:"required"`
		Campaign    Campaign    `json:"campaign" validate:"required"`
		Offset      civi.Offset `json:"offset" validate:"required"`
		RedirectURL string      `json:"redirect_url" validate:"required"`
	}
	//NewToken represents the user request for get a new token by mail
	NewToken struct {
		Campaign Campaign `json:"campaign" validate:"required"`
		Email    string   `json:"email" validate:"required"`
	}
	//SignUpReturn represents the return value for SignUp dao.
	SignUpReturn struct {
		User  vmod.User
		Token vmod.Token
	}
)

//CrmUser creates CrmUserSignUp model from SignUp
func (s *SignUp) CrmUser(dropsID string, link string) *civi.CrmUserSignUp {
	crmUser := new(civi.CrmUserSignUp)
	crmUser.CrmData.DropsId = dropsID
	crmUser.CrmData.CampaignId = s.Campaign.CampaignID
	crmUser.CrmData.Activity = "DROPS_USER_CREATE"
	crmUser.CrmData.Created = time.Now().Unix()
	crmUser.CrmUser.Email = s.SignUser.Email
	crmUser.CrmUser.FirstName = s.SignUser.FirstName
	crmUser.CrmUser.LastName = s.SignUser.LastName
	crmUser.CrmUser.PrivacyPolicy = s.SignUser.PrivacyPolicy
	crmUser.CrmUser.Country = s.SignUser.Country
	crmUser.Mail.Link = link
	crmUser.Offset = s.Offset
	return crmUser

}

//CrmEmailBody creates CrmEmailBody for handling password reset, confirmation mail and thinks like that.
func (nt *NewToken) CrmEmailBody(a string) *crm.CrmEmailBody {
	ce := new(crm.CrmEmailBody)
	ce.CrmData.CampaignId = nt.Campaign.CampaignID
	ce.CrmData.Activity = a
	ce.Mail.Email = nt.Email
	return ce
}

//User creates vmod.User from SignUser
func (u *SignUser) User(cTime int64) *vmod.User {
	//create password models based on bcrypt
	modified := vmod.InitModified(cTime)
	// create bcrypt.User
	user := new(vmod.User)
	user.ID = uuid.New().String()
	user.Email = u.Email
	user.Permission = *vmod.InitPermission("joined", "drops-service", modified.Created)
	if u.Service != "" && u.Service != "drops-service" {
		user.Permission.Add("member", u.Service, modified.Created)
	}
	user.Modified = *modified
	return user
}

//Profile creates Profle model from SignUser
func (u *SignUser) Profile(cTime int64, userID string) *vmod.Profile {
	// create Profile
	modified := vmod.InitModified(cTime)
	profile := new(vmod.Profile)
	profile.ID = uuid.New().String()
	profile.UserID = userID
	profile.FirstName = u.FirstName
	profile.LastName = u.LastName
	profile.FullName = u.FirstName + " " + u.LastName
	profile.Country = u.Country
	profile.Modified = *modified
	return profile
}

//Policies creates policies model for from SignUser. Contains `confiremd` and `privacy_policy` state.
func (u *SignUser) Policies(cTime int64, userID string) *vmod.Policies {
	policies := vmod.InitPolicies(userID, "confirmed", cTime)
	policies.Add("privacy_policy", u.PrivacyPolicy, cTime)
	return policies
}
