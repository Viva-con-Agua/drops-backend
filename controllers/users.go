package controllers

import (
	"net/http"

	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

/**
 * "GET /users"
 * response list of models.User
 */
func UserList(c echo.Context) (err error) {
	query := new(models.UserQuery)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.UserList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
* "GET /users/:uuid"
* response models.User
 */
func UserById(c echo.Context) (err error) {
	uuid := c.Param("uuid")
	response, err := database.UserById(uuid)
	if err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
 * "GET /users"
 * response uuid of updated models.User
 */
func UserUpdate(c echo.Context) (err error) {
	// create body as models.User
	body := new(models.User)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.UserUpdate(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Updated(body.Uuid))
}

/**
 * "DELETE /users"
 * response uuid of deleted models.User
 */
func UserDelete(c echo.Context) (err error) {
	// create body as models.User
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
	if err = database.UserDelete(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Deleted(body.Uuid))
}
