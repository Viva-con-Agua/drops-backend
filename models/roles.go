package models

import (
	"strconv"
)

type (
	CrewRoleList []CrewRole
	CrewRole     struct {
		Uuid        string `json:"uuid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	CrewRoleCreate struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	CrewRoleUpdate struct {
		Uuid        string `json:"uuid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	CrewRoleExtended struct {
		Uuid        string       `json:"uuid" validate:"required"`
		Name        string       `json:"name" validate:"required"`
		Description string       `json:"description" validate:"required"`
		Permissions []Permission `json:"permissions" validate:"required"`
	}
	AssignCrewRole struct {
		ProfileId string `json:"profile_id" validate:"required"`
		CrewId    string `json:"crew_id" validate:"required"`
		RoleId    string `json:"role_id" validate:"required"`
	}
	QueryRole struct {
		Offset string `query:"offset" default:"0"`
		Count  string `query:"count" default:"40"`
		Name   string `query:"name" default:"%"`
		Sort   string `query:"sort"`
		SortBy string `query:"sortby"`
	}
	FilterRole struct {
		Name string
	}
)

func (list *CrewRoleList) Distinct() *CrewRoleList {
	r := make(CrewRoleList, 0, len(*list))
	m := make(map[CrewRole]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}

func (list *CrewRoleList) NotEmpty() *CrewRoleList {
	r := make(CrewRoleList, 0, len(*list))
	for _, val := range *list {
		if val.Uuid != "" {
			r = append(r, val)
		}
	}
	return &r
}

func (q *QueryRole) Defaults() {
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

func (q *QueryRole) Page() *Page {
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

func (q *QueryRole) OrderBy() string {
	var asc = "ASC"
	if q.Sort == "DESC" {
		asc = " DESC"
	}
	var sort = "ORDER BY "
	if q.SortBy == "" {
		return ""
	}
	if q.SortBy == "name" {
		return sort + " cr.name " + asc
	}
	return sort
}

func (q *QueryRole) Filter() *FilterRole {
	filter := new(FilterRole)
	if q.Name != "" {
		filter.Name = q.Name
	} else {
		filter.Name = "%"
	}
	return filter
}
