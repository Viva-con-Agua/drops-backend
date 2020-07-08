package models

import (
	//	"database/sql"
	"strconv"
)

type (
	ProfileCreate struct {
		Email     string `json:"email" validate:"required"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Mobile    string `json:"birthdate" validate:"required"`
		Birthdate int    `json:"birthdate" validate:"required"`
		Gender    string `json:"gender" validate:"required"`
	}
	ProfileUpdate struct {
		Uuid      string `json:"uuid" validate:"required"`
		Email     string `json:"email" validate:"required"`
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Mobile    string `json:"mobilePhone" validate:"required"`
		Birthdate int    `json:"birthday" validate:"required"`
		Gender    string `json:"gender" validate:"required"`
	}

	/*
	 * Models for profiles
	 */
	ProfileDefault struct {
		Uuid      string    `json:"uuid" validate:"required"`
		Email     string    `json:"email" validate:"required"`
		FirstName string    `json:"first_name" validate:"required"`
		LastName  string    `json:"last_name" validate:"required"`
		FullName  string    `json:"full_name" validate:"required"`
		Mobile    string    `json:"mobile_phone" validate:"required"`
		Birthdate int       `json:"birthdate" validate:"required"`
		Gender    string    `json:"gender" validate:"required"`
		Addresses []Address `json:"addresses" validate:"required"`
		Updated   int       `json:"updated" validate:"required"`
		Created   int       `json:"created" validate:"required"`
	}
	ProfileId struct {
		Uuid string `json:"uuid" validate:"required"`
	}
	ProfileMin struct {
		Uuid      string `json:"uuid" validate:"required"`
		Email     string `json:"email" validate:"required"`
		FullName  string `json:"full_name" validate:"required"`
		Mobile    string `json:"mobile_phone" validate:"required"`
		Birthdate int    `json:"birthdate" validate:"required"`
		Gender    string `json:"gender" validate:"required"`
		//PrimaryCrew sql.NullString `json:"crew" validate:"required"`
		PrimaryCrew string `json:"crew" validate:"required"`
		Created     int    `json:"created" validate:"required"`
		Avatar      Avatar `json:"avatar" validate:"required"`
	}
	ProfileExtended struct {
		Uuid      string     `json:"uuid" validate:"required"`
		Email     string     `json:"email" validate:"required"`
		Avatar    Avatar     `json:"avatar" validate:"required"`
		FirstName string     `json:"first_name" validate:"required"`
		LastName  string     `json:"last_name" validate:"required"`
		FullName  string     `json:"full_name" validate:"required"`
		Mobile    string     `json:"mobile_phone" validate:"required"`
		Birthdate int        `json:"birthdate" validate:"required"`
		Gender    string     `json:"gender" validate:"required"`
		Crews     []Crew     `json:"crews" validate:"required"`
		Addresses []Address  `json:"addresses" validate:"required"`
		Roles     []CrewRole `json:"roles" validate:"required"`
		Updated   int        `json:"updated" validate:"required"`
		Created   int        `json:"created" validate:"required"`
	}
	Avatar struct {
		Url  string `json:"url" validate:"required"`
		Type string `json:"type" validate:"required"`
	}
	QueryProfile struct {
		Offset string `query:"offset" default:"0"`
		Count  string `query:"count" default:"40"`
		Email  string `query:"email" default:"%"`
		// TODO FILTER FOR NAME AND AGE/SEX?
		Sort   string `query:"sort"`
		SortBy string `query:"sortby"`
	}
	FilterProfile struct {
		Email string
	}
)

func (q *QueryProfile) Defaults() {
	if q.Offset == "" {
		q.Offset = "0"
	}
	if q.Count == "" {
		q.Count = "40"
	}
	if q.Email == "" {
		q.Email = "%"
	}
}

func (q *QueryProfile) Page() *Page {
	var err error
	page := new(Page)
	page.Offset, err = strconv.Atoi(q.Offset)
	if err != nil {
		page.Offset = 0
	}
	page.Count, err = strconv.Atoi(q.Count)
	if err != nil {
		page.Count = 40
	}
	return page
}

func (q *QueryProfile) OrderBy() string {
	var asc = "ASC"
	if q.Sort == "DESC" {
		asc = " DESC"
	}
	var sort = "ORDER BY "
	if q.SortBy == "" {
		return ""
	}
	if q.SortBy == "email" {
		return sort + " p.email " + asc
	}
	return sort
}

func (q *QueryProfile) Filter() *FilterProfile {
	filter := new(FilterProfile)
	if q.Email != "" {
		filter.Email = q.Email
	} else {
		filter.Email = "%"
	}
	return filter
}
