package controllers

import (
	"log"
	"net/http"

	"../database"
	"../models"
	"github.com/labstack/echo"
)

type ()

/**
 * join user to role
 */
func JoinUserRole(c echo.Context) (err error) {

	// create body as models.Role
	body := new(models.UserRole)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// insert body into database
	if err = database.JoinUserRole(body); err != nil {
		return c.JSON(http.StatusInternalServerError, models.InternelServerError)
	}
	// response created
	return c.JSON(http.StatusCreated, models.Created())
}

/**
 * Response list of models.User
 */
func GetUserList(c echo.Context) (err error) {
	query := new(models.QueryUser)

	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	log.Print(filter)
	response, err := database.GetUserList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}