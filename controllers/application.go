package controllers

import (
	"drops-backend/dao"
	"drops-backend/models"
	"net/http"

	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/labstack/echo/v4"
)

// CreateApplication web controller that handles a application creation.
func CreateApplication(c echo.Context) (err error) {
	body := new(models.ApplicationCreate)
	if apiErr := verr.JSONValidate(c, body); apiErr != nil {
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	resp, apiErr := dao.CreateApplication(*body)
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	return c.JSON(http.StatusCreated, resp)
}

//GetApplicationList provides read and filter application lists
func GetApplicationList(c echo.Context) (err error) {
	resp, apiErr := dao.GetApplicationList()
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	return c.JSON(http.StatusOK, resp)
}
