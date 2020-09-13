package models

import (
	"github.com/Viva-con-Agua/echo-pool/api"
)

type (
	User struct {
		Uuid       string         `json:"uuid"`
		Email      string         `json:"email"`
		Confirmed  int            `json:"confirmed"`
		Access     api.AccessList `json:"access"`
		Profile    Profile        `json:"profile"`
		Updated    int64          `json:"updated"`
		Created    int64          `json:"created"`
		Additional api.Additional `json:"additional"`
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
	UserListFilter struct {
		UserList []string `json:"user_list" validate:"required"`
	}
)

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
