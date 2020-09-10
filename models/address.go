package models

import (
	"strconv"
)

type (
	QueryAddress struct {
		Offset   string `query:"offset" default:"0"`
		Count    string `query:"count" default:"40"`
		Street   string `query:"street" default:"%"`
		Zip      string `query:"zip" default:"%"`
		City     string `query:"city" default:"%"`
		Country  string `query:"country" default:"%"`
		GoogleId string `query:"google_id" default:"%"`
		// TODO FILTER FOR NAME AND AGE/SEX?
		Sort   string `query:"sort"`
		SortBy string `query:"sortby"`
	}
	FilterAddress struct {
		Street string
	}
)

func (q *QueryAddress) Defaults() {
	if q.Offset == "" {
		q.Offset = "0"
	}
	if q.Count == "" {
		q.Count = "40"
	}
	if q.Street == "" {
		q.Street = "%"
	}
	if q.Zip == "" {
		q.Zip = "%"
	}
	if q.City == "" {
		q.City = "%"
	}
	if q.Country == "" {
		q.Country = "%"
	}
	if q.GoogleId == "" {
		q.GoogleId = "%"
	}
}

func (q *QueryAddress) Page() *Page {
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

func (q *QueryAddress) OrderBy() string {
	var asc = "ASC"
	if q.Sort == "DESC" {
		asc = " DESC"
	}
	var sort = "ORDER BY "
	if q.SortBy == "" {
		return ""
	}
	if q.SortBy == "street" {
		return sort + " a.street " + asc
	}
	if q.SortBy == "zip" {
		return sort + " a.zip " + asc
	}
	if q.SortBy == "city" {
		return sort + " a.city " + asc
	}
	if q.SortBy == "country" {
		return sort + " a.country " + asc
	}
	return sort
}

func (q *QueryAddress) Filter() *FilterAddress {
	filter := new(FilterAddress)
	if q.Street != "" {
		filter.Street = q.Street
	} else {
		filter.Street = "%"
	}
	return filter
}
