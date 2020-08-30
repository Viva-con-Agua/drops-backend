package controllers

import (
	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/nats"
	"drops-backend/utils"
	"log"
	"net/http"

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
	user_uuid, access_token, err := database.SignUp(body)
	if err != nil {
		if err == utils.ErrorConflict {
			return c.JSON(http.StatusConflict, resp.Conflict())
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	//signin user
	user, err := database.GetSessionUser(user_uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	//TODO iRobert Request CrmUser

	log.Print(body.CrmUserSignUp(*user_uuid, *access_token))

	auth.SetSession(c, user)
	// response created
	return c.JSON(http.StatusCreated, user)
}

func ConfirmSignUp(c echo.Context) (err error) {
	token := c.Param("token")
	user_uuid, err := database.ConfirmSignUp(token)
	if err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNotFound, resp.NoContent(token))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	user, err := database.GetSessionUser(user_uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	auth.SetSession(c, user)

	//TODO iRobert Request   Activity: Confirmed Account
	return c.JSON(http.StatusOK, user)
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
	body := new(models.SignInData)
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
	auth.SetSession(c, user)
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
	return c.JSON(http.StatusOK, user)
}

func SignOut(c echo.Context) (err error) {
	auth.DeleteSession(c)
	msg := new(resp.ResponseMessage)
	msg.Message = "Successful sign out"
	return c.JSON(http.StatusOK, msg)
}
