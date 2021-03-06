package controllers

import (
	"drops-backend/database"
	"drops-backend/models"
	"net/http"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo/v4"
)

func ModelCreate(c echo.Context) (err error) {
	body := new(api.ModelCreate)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	// insert body into database
	model, api_err := database.ModelCreate(body)
	if api_err.Error != nil {
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	// response created
	return c.JSON(http.StatusCreated, model)
}

func ModelDelete(c echo.Context) (err error) {
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
	if err = database.ModelDelete(body); err != nil {
		if err == api.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Deleted(body.Uuid))
}
