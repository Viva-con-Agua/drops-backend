package controllers

import (
	"../database"
	"../models"
	"../utils"
	"errors"
	"github.com/Viva-con-Agua/echo-pool/pool"
	"github.com/labstack/echo"
	"net/http"
)

/**
 * Response list of models.Crew
 */
func GetCrewList(c echo.Context) (err error) {
	query := new(models.QueryCrew)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	query.Defaults()
	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.GetCrewList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func CreateCrew(c echo.Context) (err error) {
	// create body as models.ProfileCreate
	body := new(models.CrewCreate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.CreateCrew(body); err != nil {
		if err == utils.ErrorConflict {
			return c.JSON(http.StatusNoContent, pool.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Created())
}

func ReadCrew(c echo.Context) (err error) {
	uuid := c.Param("id")
	response, err := database.GetCrew(uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
	}
	if response == nil {
		return c.JSON(http.StatusNoContent, pool.NoContent(uuid))
	}
	return c.JSON(http.StatusOK, response)
}

func UpdateCrew(c echo.Context) (err error) {
	// create body as models.Crew
	body := new(models.CrewUpdate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.UpdateCrew(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Updated(body.Uuid))
}

func DeleteCrew(c echo.Context) (err error) {
	// create body as models.Crew
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
	if err = database.DeleteCrew(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Deleted(body.Uuid))
}

func AssignCrew(c echo.Context) (err error) {
	// create body as models.ProfileCreate
	body := new(models.AssignCrew)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.AssignCrew(body); err != nil {
		if err == utils.ErrorConflict {
			return c.JSON(http.StatusNoContent, pool.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Created())
}

func RemoveCrew(c echo.Context) (err error) {
	// create body as models.ProfileCreate
	body := new(models.RemoveCrew)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.RemoveCrew(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.ProfileId+"_"+body.CrewId))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Deleted(body.ProfileId+"_"+body.CrewId))
}

/**
 * Response list of models.CrewRoleList
 */
func ActiveStateChange(c echo.Context) (err error) {
	// create body as models.ProfileProfileNewsletter
	body := new(models.ActiveState)
	state := c.Param("state")

	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if state != "requested" && state != "active" && state != "inactive" {
		return c.JSON(http.StatusBadRequest, errors.New("Invalid state"))
	}

	// update body into database
	if err = database.ActiveStateChange(body, state); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.ProfileId))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Updated(body.ProfileId))
}

/**
 * Non voting membership
 */
func NVMStateChange(c echo.Context) (err error) {
	// create body as models.ProfileProfileNewsletter
	body := new(models.NVMState)
	state := c.Param("state")

	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if state != "active" && state != "inactive" {
		return c.JSON(http.StatusBadRequest, errors.New("Invalid state"))
	}

	// update body into database
	if err = database.NVMStateChange(body, state); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, pool.NoContent(body.ProfileId))
		}
		return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, pool.Updated(body.ProfileId))
}
