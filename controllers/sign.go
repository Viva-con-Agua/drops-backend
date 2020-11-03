package controllers

import (
	"drops-backend/dao"
	"drops-backend/models"
	"net/http"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

func SignUp(c echo.Context) (err error) {
	body := new(models.SignUp)
	if resp := vcago.JsonErrorHandling(c, body); resp != nil {
		return c.JSON(resp.Code, resp.Response)
	}
	// insert body into database
	user, api_err := dao.SignUp(body)
	if resp := vcago.ResponseErrorHandling(c, api_err, body); resp != nil {
		return c.JSON(resp.Code, resp.Response)
	}
	//signin user
	return c.JSON(http.StatusCreated, user)
}

func SignIn(c echo.Context) (err error) {
	body := new(models.SignIn)
	if resp := vcago.JsonErrorHandling(c, body); resp != nil {
		return c.JSON(resp.Code, resp.Response)
	}
	user, api_err := dao.SignIn(body)
	if resp := vcago.ResponseErrorHandling(c, api_err, body); resp != nil {
		return c.JSON(resp.Code, resp.Response)
	}
	return c.JSON(http.StatusOK, user)
}
