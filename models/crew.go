package models

import (
	"strconv"
)

type (
	CrewCreate struct {
		Name         string `json:"name" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Abbreviation string `json:"abbreviation" validate:"required"`
		Cities       []City `json:"cities" validate:"required"`
	}

	CrewUpdate struct {
		Uuid         string `json:"uuid" validate:"required"`
		Name         string `json:"name" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Abbreviation string `json:"abbreviation" validate:"required"`
		Cities       []City `json:"cities"`
	}

	Crew struct {
		Uuid         string `json:"uuid" validate:"required"`
		Primary      int    `json:"primary" validate:"required"`
		Name         string `json:"name" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Abbreviation string `json:"abbreviation" validate:"required"`
	}
	CrewExtended struct {
		Uuid         string `json:"uuid" validate:"required"`
		Name         string `json:"name" validate:"required"`
		Email        string `json:"email" validate:"required"`
		Abbreviation string `json:"abbreviation" validate:"required"`
		Cities       []City `json:"cities"`
		CrewMeta     []Dict `json:"meta"`
		Updated      int    `json:"updated" validate:"required"`
		Created      int    `json:"created" validate:"required"`
	}
	CrewFull struct {
		Uuid         string            `json:"uuid" validate:"required"`
		Name         string            `json:"name" validate:"required"`
		Email        string            `json:"email" validate:"required"`
		Abbreviation string            `json:"abbreviation" validate:"required"`
		Cities       []City            `json:"cities"`
		CrewMeta     []Dict            `json:"meta"`
		Roles        []CrewRoleProfile `json:"roles"`
		Updated      int               `json:"updated" validate:"required"`
		Created      int               `json:"created" validate:"required"`
	}
	CrewRoleProfile struct {
		Id          int    `json:"id" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
	}
	CrewRoleProfileList []CrewRoleProfile
	AssignCrew          struct {
		ProfileId string `json:"profile_id" validate:"required"`
		CrewId    string `json:"crew_id" validate:"required"`
		Primary   bool   `json:"primary" validate:"required"`
		Active    bool   `json:"active" validate:"required"`
	}
	ActiveState struct {
		ProfileId string `json:"profile_id" validate:"required"`
		CrewId    string `json:"crew_id" validate:"required"`
	}
	NVMState struct {
		ProfileId string `json:"profile_id" validate:"required"`
		CrewId    string `json:"crew_id" validate:"required"`
	}
	RemoveCrew struct {
		ProfileId string `json:"profile_id" validate:"required"`
		CrewId    string `json:"crew_id" validate:"required"`
	}

	QueryCrew struct {
		Offset string `query:"offset" default:"0"`
		Count  string `query:"count" default:"40"`
		Name   string `query:"name" default:"%"`
		Sort   string `query:"sort"`
		SortBy string `query:"sortby"`
	}

	CrewList         []Crew
	CrewExtendedList []CrewExtended

	FilterCrew struct {
		Name string
	}
)

func (list *CrewList) Distinct() *CrewList {
	r := make(CrewList, 0, len(*list))
	m := make(map[Crew]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}

func (q *QueryCrew) Defaults() {
	if q.Offset == "" {
		q.Offset = "0"
	}
	if q.Count == "" {
		q.Count = "40"
	}
	if q.Name == "" {
		q.Name = "%"
	}
}

func (q *QueryCrew) Page() *Page {
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

func (q *QueryCrew) OrderBy() string {
	var asc = "ASC"
	if q.Sort == "DESC" {
		asc = " DESC"
	}
	var sort = "ORDER BY "
	if q.SortBy == "" {
		return ""
	}
	if q.SortBy == "name" {
		return sort + " p.name " + asc
	}
	return sort
}

func (q *QueryCrew) Filter() *FilterCrew {
	filter := new(FilterCrew)
	if q.Name != "" {
		filter.Name = q.Name
	} else {
		filter.Name = "%"
	}
	return filter
}
