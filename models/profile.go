package models

import (
	//	"database/sql"

	"strconv"
	"strings"

	"github.com/Viva-con-Agua/vcago/vmod"
)

type (
	Avatar struct {
		Url     string `json:"url"`
		Type    string `json:"type"`
		Updated int64  `json:"updated" validate:"required"`
		Created int64  `json:"created" validate:"required"`
	}
	ProfileCreate struct {
		Uuid      string `json:"profile_uuid" validate:"required"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Updated   int64  `json:"updated" validate:"required"`
		Created   int64  `json:"created" validate:"required"`
	}
	Profile struct {
		Uuid        string `json:"uuid" validate:"required"`
		UserUuid    string `json:"user_uuid" validate:"required"`
		Avatar      Avatar `json:"avatar"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
		FullName    string `json:"full_name" validate:"required"`
		DisplayName string `json:"display_name" validate:"required"`
		Gender      string `json:"gender"`
		Updated     int64  `json:"updated" validate:"required"`
		Created     int64  `json:"created" validate:"required"`
	}
	QueryProfile struct {
		Offset string `query:"offset" default:"0"`
		Count  string `query:"count" default:"40"`
		Email  string `query:"email" default:"%"`
		Sort   string `query:"sort"`
		SortBy string `query:"sortby"`
	}
	FilterProfile struct {
		Email string
	}
	ListRequest struct {
		UuidList []string `json:"uuid_list" validate:"required"`
	}
	ProfileDB struct {
		ID          string        `bson:"_id" json:"profile_id" validate:"required"`
		UserID      string        `bson:"user_id" json:"user_id" validate:"required"`
		FirstName   string        `bson:"first_name" json:"first_name" validate:"required"`
		LastName    string        `bson:"last_name" json:"last_name" validate:"required"`
		FullName    string        `bson:"full_name" json:"full_name" validate:"required"`
		DisplayName string        `bson:"display_name" json:"display_name" validate:"required"`
		Gender      string        `bson:"gender" json:"gender"`
		Country     string        `bson:"country" json:"country"`
		Avatar      vmod.Avatar   `bson:"avatar" json:"avatar"`
		Modified    vmod.Modified `json:"modified" bson:"modified" validation:"required"`
	}
)

func (p *ProfileDB) Profile() *vmod.Profile {
	profile := new(vmod.Profile)
	profile.ID = p.ID
	profile.UserID = p.UserID
	profile.FirstName = p.FirstName
	profile.LastName = p.LastName
	profile.FullName = p.FullName
	profile.DisplayName = p.DisplayName
	profile.Gender = p.Gender
	profile.Country = p.Country
	profile.Avatar = p.Avatar
	profile.Modified = p.Modified
	return profile
}

func (l *ListRequest) Filter() string {
	if l.UuidList != nil {
		filter := "WHERE "
		for _, val := range l.UuidList {
			filter = filter + "du.uuid = '" + val + "' OR "
		}
		filter = strings.TrimSuffix(filter, "OR ")
		return filter
	} else {
		return ""
	}
}

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
