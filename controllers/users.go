package controllers

import (
	"log"
	"net/http"
	"os"

	"drops-backend/crm"
	"drops-backend/database"
	"drops-backend/models"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo/v4"
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

func PasswordResetToken(c echo.Context) (err error) {
	body := new(models.NewToken)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	crm_email, api_err := database.PasswordResetToken(body)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorUserNotFound {
			return c.JSON(http.StatusNotFound, api.RespNoContent("email", body.Email))
		}
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	if os.Getenv("CRM_SIGNUP") != "false" {
		err = crm.IrobertSendMail(crm_email)
	} else {
		log.Print(crm_email)
	}
	//TODO CRM Request new Token for Signup
	return c.JSON(http.StatusCreated, api.RespCreated())
}

func PasswordReset(c echo.Context) (err error) {
	body := new(models.PasswordReset)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	u_uuid, api_err := database.PasswordReset(body)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorUserNotFound {
			return c.JSON(http.StatusNotFound, api.RespNoContent("token", body.Token))
		}
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	key := "u_uuid"
	//TODO CRM Request new Token for Signup
	return c.JSON(http.StatusOK, api.RespCustom("password_updated", &key, *u_uuid))
}
