package controllers

import (
	"../database"
	"../models"
	"../utils"
	"github.com/Viva-con-Agua/echo-pool/pool"
	"github.com/labstack/echo"
	"net/http"
)

/**
 * Response list of models.ProfileDefaultList
 */
func GetProfileDefaultList(c echo.Context) (err error) {
	query := new(models.QueryProfile)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	query.Defaults()
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.GetProfileDefaultList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
 * Response list of models.ProfileId
 */
func GetProfileIdList(c echo.Context) (err error) {
	query := new(models.QueryProfile)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	query.Defaults()
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.GetProfileIdList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
 * Response list of models.ProfileMin
 */
func GetProfileMinList(c echo.Context) (err error) {
	query := new(models.QueryProfile)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	query.Defaults()
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.GetProfileMinList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
 * Response list of models.ProfileMin
 */
func GetProfileExtendedList(c echo.Context) (err error) {
	query := new(models.QueryProfile)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	query.Defaults()
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.GetProfileExtendedList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func CreateProfile(c echo.Context) (err error) {
	// create body as models.ProfileCreate
	body := new(models.ProfileCreate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.CreateProfile(body); err != nil {
		if err == utils.ErrorConflict {
			return c.JSON(http.StatusNoContent, pool.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Created())
}

func ReadProfile(c echo.Context) (err error) {
	uuid := c.Param("id")
	response, err := database.GetProfile(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	if response == nil {
		return c.JSON(http.StatusNoContent, pool.NoContent(uuid))
	}
	return c.JSON(http.StatusOK, response)
}

func UpdateProfile(c echo.Context) (err error) {
	// create body as models.Profile
	body := new(models.ProfileUpdate)

	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.UpdateProfile(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Updated(body.Uuid))
}

func DeleteProfile(c echo.Context) (err error) {
	// create body as models.DeleteBody
	body := new(models.DeleteBody)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.DeleteProfile(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Deleted(body.Uuid))
}

/**
 * Response list of models.CrewRoleList
 */
func SetNewsletter(c echo.Context) (err error) {
	// create body as models.ProfileProfileNewsletter
	body := new(models.ProfileNewsletter)

	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.SetNewsletter(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Updated(body.Uuid))
}
