package controllers

import (
	"drops-backend/database"
	"drops-backend/models"
	"net/http"
	"strings"

	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

func ServiceList(c echo.Context) (err error) {
	response, err := database.ServiceList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func ServiceInsert(c echo.Context) (err error) {
	body := new(models.ServiceCreate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// insert body into database
	response, err := database.ServiceInsert(body)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.JSON(http.StatusConflict, resp.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusCreated, response)
}
