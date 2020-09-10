package controllers

import (
	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/utils"
	"net/http"

	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

func AccessInsert(c echo.Context) (err error) {
	// create body as models.Role
	body := new(models.AccessCreate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// insert body into database
	if err = database.AccessInsert(body); err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	// response created
	return c.JSON(http.StatusCreated, resp.Created())
}

func AccessDelete(c echo.Context) (err error) {
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
	if err = database.AccessDelete(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Deleted(body.Uuid))
}
