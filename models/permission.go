package models

import (
	"net/http"

	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type (
	PermissionList []Permission
	Permission     struct {
		Uuid        string `json:"uuid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
)

func (list *PermissionList) Distinct() *PermissionList {
	r := make(PermissionList, 0, len(*list))
	m := make(map[Permission]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}

func (list *PermissionList) NotEmpty() *PermissionList {
	r := make(PermissionList, 0, len(*list))
	for _, val := range *list {
		if val.Uuid != "" {
			r = append(r, val)
		}
	}
	return &r
}

func (r *PermissionList) Permission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		if sess.Values["roles"] == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, resp.Unauthorized())
		}
		return next(c)
	}
}
