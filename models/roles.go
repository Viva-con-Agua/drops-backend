package models

import (
	"github.com/Viva-con-Agua/echo-pool/pool"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

type (
	CrewRoleList []CrewRole

	CrewRole struct {
		Uuid string `json:"crew_uuid" validate:"required"`
		Crew string `json:"crew_name" validate:"required"`
		Role string `json:"role" validate:"required"`
	}

	RoleCreate struct {
		Name   string `json:"name" validate:"required"`
		Pillar string `json:"pillar" validate:"required"`
	}
	Role struct {
		Uuid   string `json:"uuid" validate:"required"`
		Name   string `json:"name" validate:"required"`
		Pillar string `json:"pillar" validate:"required"`
	}
	Roles struct {
		Role []Role
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
		if val.Role != "" {
			r = append(r, val)
		}
	}
	return &r
}

func (roles *Roles) AddRole(role Role) []Role {
	roles.Role = append(roles.Role, role)
	return roles.Role
}

func (r *Roles) Permission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		if sess.Values["roles"] == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, pool.Unauthorized())
		}
		return next(c)
	}
}
