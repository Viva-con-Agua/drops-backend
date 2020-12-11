package controllers

import (
	"drops-backend/dao"
	"drops-backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

//GetProfileList web controller for providing dao.GetProfileList
func GetProfileList(c echo.Context) (err error) {
	query := new(models.ProfileQuery)
	if err := c.Bind(query); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	resp, apiErr := dao.GetProfileList(*query)
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	return c.JSON(http.StatusOK, resp)
}
