package models

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
		Kampagne    int          `json:"kampagne" validate:"required"`
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
		Token string `json:"token" validate:"required"`
	}
	CrmUser struct {
		Uuid      string       `json:"drops_id" validate:"required"`
		Email     string       `json:"email" validate:"required"`
		FirstName string       `json:"first_name" validate:"required"`
		LastName  string       `json:"last_name" validate:"required"`
		Kampagne  int          `json:"kampagne" validate:"required"`
		Mail      Mail         `json:"mail" validate:"required"`
		Offset    SignUpOffset `json:"offset" validate:"required"`
	}
)

func (signup_data *SignUpData) CrmUser(drops_id string, token string) *CrmUser {
	crm_user := new(CrmUser)
	crm_user.Uuid = drops_id
	crm_user.Email = signup_data.Email
	crm_user.FirstName = signup_data.FirstName
	crm_user.LastName = signup_data.LastName
	crm_user.Kampagne = signup_data.Kampagne
	crm_user.Mail.Token = token
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
