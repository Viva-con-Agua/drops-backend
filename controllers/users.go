package controllers

import (
	"net/http"

	"drops-backend/database"
	"drops-backend/models"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

/**
 * "GET /users"
 * response list of models.User
 */
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
		if err == api.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Deleted(body.Uuid))
}
