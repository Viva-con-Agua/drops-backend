package controllers

import (
	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/nats"
	"drops-backend/utils"
	"log"
	"net/http"
	"strings"

	"github.com/Viva-con-Agua/echo-pool/auth"
	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

func SignUp(c echo.Context) (err error) {
	body := new(models.SignUpData)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// insert body into database
	access_token, err := database.SignUp(body)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.JSON(http.StatusConflict, resp.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	mail := models.MailInfo{To: body.Email, Token: *access_token, Template: "default"}
	log.Print(mail)
	nats.PublishToken(&mail)
	// response created
	return c.JSON(http.StatusCreated, resp.Created())
}

func ConfirmSignUp(c echo.Context) (err error) {
	token := c.Param("token")
	err = database.ConfirmSignUp(token)
	if err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNotFound, resp.NoContent(token))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	response := new(resp.ResponseMessage)
	response.Message = "Successful Confirmed"
	return c.JSON(http.StatusOK, response)
}

func SignUpToken(c echo.Context) (err error) {
	body := new(models.NewToken)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	access_token, err := database.SignUpToken(body)
	if err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNotFound, resp.NoContent(body.Email))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	mail := models.MailInfo{To: body.Email, Token: *access_token, Template: "default"}
	log.Print(mail)
	nats.PublishToken(&mail)
	// response created
	return c.JSON(http.StatusCreated, resp.Created())
}

func SignIn(c echo.Context) (err error) {
	body := new(models.Credentials)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user, err := database.SignIn(body)
	if err != nil {
		if err == utils.ErrorPassword {
			return c.JSON(http.StatusUnauthorized, &resp.ResponseMessage{Message: "No valid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	if user == nil {
		return c.JSON(http.StatusUnauthorized, &resp.ResponseMessage{Message: "No valid email or password"})
	}
	if user.Confirmed == 0 {
		return c.JSON(http.StatusForbidden, &resp.ResponseMessage{Message: "Not confirmed"})
	}
	auth.SetSession(c, user, &auth.AccessToken{AccessToken: "null"})
	return c.JSON(http.StatusOK, user)
}

func Current(c echo.Context) (err error) {
	user, _ := auth.GetUser(c)
	if err != nil {
		if err == utils.ErrorPassword {
			return c.JSON(http.StatusUnauthorized, &resp.ResponseMessage{Message: "No user sign in"})
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	if user == nil {
		return c.JSON(http.StatusUnauthorized, &resp.ResponseMessage{Message: "No user sign in"})
	}
	if user.Confirmed == 0 {
		return c.JSON(http.StatusForbidden, &resp.ResponseMessage{Message: "Not confirmed"})
	}
	return c.JSON(http.StatusOK, user)
}

func SignOut(c echo.Context) (err error) {
	auth.DeleteSession(c)
	msg := new(resp.ResponseMessage)
	msg.Message = "Successful sign out"
	return c.JSON(http.StatusOK, msg)
}
