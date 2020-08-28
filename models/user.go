package models

import "github.com/Viva-con-Agua/echo-pool/crm"

type (
	Credentials struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	SignUpData struct {
		Email       string       `json:"email" validate:"required"`
		Password    string       `json:"password" validate:"required"`
		FirstName   string       `json:"first_name" validate:"required"`
		LastName    string       `json:"last_name" validate:"required"`
		ServiceName string       `json:"service_name" validate:"required"`
		RedirectUrl string       `json:"redirect_url" validate:"required"`
		Campaign    int          `json:"campaign" validate:"required"`
		Offset      SignUpOffset `json:"offset" validate:"required"`
	}

	NewToken struct {
		Email string `json:"email" validate:"required"`
	}

	User struct {
		Uuid      string   `json:"uuid" validate:"required"`
		Email     string   `json:"email" validate:"required"`
		FirstName string   `json:"first_name" validate:"required"`
		Access    []Access `json:"access"`
		Profile   Profile  `json:"profile" validate:"required"`
		Confirmed int      `json:"confirmed"`
		Updated   int      `json:"updated" validate:"required"`
		Created   int      `json:"created" validate:"required"`
	}

	UserQuery struct {
		Offset      int    `query:"offset" default:"0"`
		Count       int    `query:"count" default:"40"`
		SortDir     string `query:"sortdir"`
		SortBy      string `query:"sortby"`
		Email       string `query:"email" default:"%"`
		UpdatedFrom int    `query:"updated_from"`
		UpdatedTo   int    `query:"updated_to"`
	}

	UserFilter struct {
		Email string
	}
	SignUpOffset struct {
		KnownFrom string `json:"known_from" validate:"required"`
	}
	Mail struct {
		Link string `json:"link" validate:"required"`
	}
	CrmUser struct {
		Email     string `json:"email" validate:"required"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}
	CrmUserSignUp struct {
		CrmData crm.CrmData  `json:"crm_data" validate:"required"`
		CrmUser CrmUser      `json:"crm_user" validate:"required"`
		Mail    Mail         `json:"mail" validate:"required"`
		Offset  SignUpOffset `json:"offset" validate:"required"`
	}
)

func (signup_data *SignUpData) CrmUserSignUp(drops_id string, link string) *CrmUserSignUp {
	crm_user := new(CrmUserSignUp)
	crm_user.CrmData.DropsId = drops_id
	crm_user.CrmData.CampaignId = signup_data.Campaign
	crm_user.CrmData.Activity = "SIGNUP"
	crm_user.CrmUser.Email = signup_data.Email
	crm_user.CrmUser.FirstName = signup_data.FirstName
	crm_user.CrmUser.LastName = signup_data.LastName
	crm_user.Mail.Link = link
	crm_user.Offset = signup_data.Offset
	return crm_user

}
func (q *UserQuery) Page() *Page {
	//create new Page
	page := new(Page)
	//set offset, default null
	page.Offset = q.Offset
	//set count, default 20
	if q.Count == 0 {
		page.Count = 20
	} else {
		page.Count = q.Count
	}
	//return Page
	return page
}

func (q *UserQuery) OrderBy() string {
	// get order direction
	var dir string
	if q.SortDir == "DESC" || q.SortDir == "ASC" {
		dir = q.SortDir + " "
	} else {
		dir = "DESC "
	}
	// return sort string
	if q.SortBy == "id" {
		return "u.id " + dir
	} else {
		return "u.id " + dir
	}
}

func (q *UserQuery) Filter() *UserFilter {
	filter := new(UserFilter)
	if q.Email != "" {
		filter.Email = q.Email
	} else {
		filter.Email = "%"
	}
	return filter
}
