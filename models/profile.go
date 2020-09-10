package models

import (
	//	"database/sql"
	"strconv"
)

type (
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
