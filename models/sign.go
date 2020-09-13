package models

import (
	"github.com/Viva-con-Agua/echo-pool/crm"
)

type (
	// for /auth/signin
	SignIn struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	SignUser struct {
		Email         string `json:"email" validate:"required"`
		Password      string `json:"password" validate:"required"`
		FirstName     string `json:"first_name" validate:"required"`
		LastName      string `json:"last_name" validate:"required"`
		PrivacyPolicy int    `json:"privacy_policy"`
	}
	Offset struct {
		KnownFrom  string `json:"known_from" validate:"required"`
		Newsletter bool   `json:"newsletter"`
	}
	Campaign struct {
		CampaignId   int    `json:"campaign_id" validate:"required"`
		CampaignName string `json:"campaign_name" validate:"required"`
	}
	SignUp struct {
		SignUser    SignUser `json:"sign_user" validate:"required"`
		Campaign    Campaign `json:"campaign" validate:"required"`
		Offset      Offset   `json:"offset" validate:"required"`
		RedirectUrl string   `json:"redirect_url" validate:"required"`
	}
	NewToken struct {
		Email string `json:"email" validate:"required"`
	}

	Mail struct {
		Email string `json:"email" validate:"required"`
		Link  string `json:"link" validate:"required"`
	}
	CrmUser struct {
		Email     string `json:"email" validate:"required"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}
	CrmUserSignUp struct {
		CrmData crm.CrmData `json:"crm_data" validate:"required"`
		CrmUser CrmUser     `json:"crm_user" validate:"required"`
		Mail    Mail        `json:"mail" validate:"required"`
		Offset  Offset      `json:"offset" validate:"required"`
	}
)

func (signup_data *SignUp) CrmUserSignUp(drops_id string, link string) *CrmUserSignUp {
	crm_user := new(CrmUserSignUp)
	crm_user.CrmData.DropsId = drops_id
	crm_user.CrmData.CampaignId = signup_data.Campaign.CampaignId
	crm_user.CrmData.Activity = "SIGNUP"
	crm_user.CrmUser.Email = signup_data.SignUser.Email
	crm_user.CrmUser.FirstName = signup_data.SignUser.FirstName
	crm_user.CrmUser.LastName = signup_data.SignUser.LastName
	crm_user.Mail.Link = link
	crm_user.Offset = signup_data.Offset
	return crm_user

}
