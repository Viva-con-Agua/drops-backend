package controllers

import (
	"drops-backend/crm"
	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/labstack/echo"
)

func SignUp(c echo.Context) (err error) {
	body := new(models.SignUp)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	// insert body into database
	user_uuid, _, api_err := database.SignUp(body)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorConflict {
			return c.JSON(http.StatusConflict, api.RespConflict("email", body.SignUser.Email))
		}
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	//signin user
	user, err_api := database.GetSessionUser(user_uuid)
	if err_api.Error != nil {
		err_api.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	api.SetSession(c, user)
	return c.JSON(http.StatusCreated, user)
}

func ConfirmSignUp(c echo.Context) (err error) {
	token := c.Param("token")
	user_uuid, api_err := database.ConfirmSignUp(token)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorNotFound {
			return c.JSON(http.StatusBadRequest, api.RespNoContent("token", token))
		}
		if strings.Contains(api_err.Error.Error(), "no rows in result set") {
			return c.JSON(http.StatusBadRequest, api.RespNoContent("token", token))
		}
		api_err.LogError(c, nil)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	user, err_api := database.GetSessionUser(user_uuid)
	if err_api.Error != nil {
		err_api.LogError(c, nil)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	api.SetSession(c, user)

	//TODO iRobert Request   Activity: Confirmed Account

	return c.JSON(http.StatusOK, user)
}

func SignUpToken(c echo.Context) (err error) {
	body := new(models.NewToken)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	access_token, api_err := database.SignUpToken(body)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorNotFound {
			return c.JSON(http.StatusNotFound, api.RespNoContent("email", body.Email))
		}
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())
	}
	err = crm.IrobertSendMail(access_token)
	log.Print(*access_token)
	//TODO CRM Request new Token for Signup
	return c.JSON(http.StatusCreated, api.RespCreated())
}

func SignIn(c echo.Context) (err error) {
	body := new(models.SignIn)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, api.JsonErrorResponse(err))
	}
	user, api_err := database.SignIn(body)
	if api_err.Error != nil {
		if api_err.Error == api.ErrorPassword {
			return c.JSON(http.StatusUnauthorized, api.RespCustom("password_false", nil, nil))
		}
		if api_err.Error == utils.ErrorUserNotFound {
			key := "email"
			return c.JSON(http.StatusUnauthorized, api.RespCustom("email_false", &key, body.Email))
		}
		if api_err.Error == utils.ErrorUserNotConfirmed {
			return c.JSON(http.StatusUnauthorized, api.RespCustom("confirmed_false", nil, nil))
		}
		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, api.RespInternelServerError())

	}
	if user.Confirmed == false {
		return c.JSON(http.StatusForbidden, api.RespCustom("Not confirmed", nil, nil))
	}
	api.SetSession(c, user)
	return c.JSON(http.StatusOK, user)
}

func Current(c echo.Context) (err error) {
	user, _ := api.GetUser(c)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, api.RespCustom("No user sign in", nil, nil))
	}
	return c.JSON(http.StatusOK, user)
}

func SignOut(c echo.Context) (err error) {
	api.DeleteSession(c)
	return c.JSON(http.StatusOK, api.RespCustom("Successful sign out", nil, nil))
}
