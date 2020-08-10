package models

import (
	"strconv"
)

type (
	AddressCreate struct {
		Street     string `json:"street" validate:"required"`
		Primary    string `json:"primary" validate;"required"`
		Additional string `json:"additional" validate;"required"`
		Zip        string `json:"zip" validate;"required"`
		City       string `json:"city" validate;"required"`
		ProfileId  string `json:"profile_id" validate;"required"`
		Country    string `json:"country" validate;"required"`
		GoogleId   string `json:"google_id" validate;"required"`
	}
	AddressUpdate struct {
		Uuid       string `json:"uuid" validate;"required"`
		Primary    string `json:"primary" validate;"required"`
		Street     string `json:"street" validate:"required"`
		Additional string `json:"additional" validate;"required"`
		Zip        string `json:"zip" validate;"required"`
		City       string `json:"city" validate;"required"`
		Country    string `json:"country" validate;"required"`
		GoogleId   string `json:"google_id" validate;"required"`
	}

	Address struct {
		Uuid       string `json:"uuid" validate:"required"`
		Primary    int    `json:"primary" validate:"required"`
		Street     string `json:"street" validate:"required"`
		Additional string `json:"additional" validate:"required"`
		Zip        string `json:"zip" validate:"required"`
		City       string `json:"city" validate:"required"`
		Country    string `json:"country" validate:"required"`
		GoogleId   string `json:"google_id" validate:"required"`
		Updated    int    `json:"updated" validate:"required"`
		Created    int    `json:"created" validate:"required"`
	}

	AddressExtended struct {
		Uuid       string         `json:"uuid" validate:"required"`
		Primary    int            `json:"primary" validate:"required"`
		Street     string         `json:"street" validate:"required"`
		Additional string         `json:"additional" validate:"required"`
		Zip        string         `json:"zip" validate:"required"`
		City       string         `json:"city" validate:"required"`
		Country    string         `json:"country" validate:"required"`
		GoogleId   string         `json:"google_id" validate:"required"`
		Profile    AddressProfile `json:"profile" validate:"required"`
		Updated    int            `json:"updated" validate:"required"`
		Created    int            `json:"created" validate:"required"`
	}
	AddressProfile struct {
		Uuid     string `json:"uuid" validate:"required"`
		Fullname string `json:"full_name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Mobile   string `json:"mobile" validate:"required"`
	}
	AddressList []Address

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

func (list *AddressList) Distinct() *AddressList {
	r := make(AddressList, 0, len(*list))
	m := make(map[Address]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
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
