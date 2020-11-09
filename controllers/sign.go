package controllers

import (
	"drops-backend/dao"
	"drops-backend/models"
	"drops-backend/nats"
	"net/http"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/verr"
	"github.com/labstack/echo/v4"
)

//SignUp is a web controller that manages the sign up process.
//Initialize the controller with the echo.Echo.PUSH() function.
func SignUp(c echo.Context) (err error) {
	body := new(models.SignUp)
	if apiErr := verr.JSONValidate(c, body); apiErr != nil {
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	// insert body into database
	resp, apiErr := dao.SignUp(body)
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	//signin user
	vcago.NewSession(c, &resp.User)
	nats.Nats.Publish("drops.signup", body.CrmUser(resp.User.ID, resp.Token.Code))
	return c.JSON(http.StatusCreated, resp.User)
}

//SignIn is a web controller that manages the sign in process.
//Initialize the controller with the echo.Echo.PUSH()
func SignIn(c echo.Context) (err error) {
	body := new(models.SignIn)
	if apiErr := verr.JSONValidate(c, body); apiErr != nil {
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	user, apiErr := dao.SignIn(body)
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	return c.JSON(http.StatusOK, user)
}

//Current return the user for current session
func Current(c echo.Context) (err error) {
	apiErr := vcago.GetSession(c)
	if apiErr != nil {
		apiErr.Log(c)
		return c.JSON(apiErr.Code, apiErr.Body)
	}
	return c.JSON(http.StatusOK, c.Get("user"))
}

//SignOut for logout the user by delete his session
func SignOut(c echo.Context) (err error) {
	vcago.DeleteSession(c)
	return c.JSON(http.StatusOK, verr.ResponseError{Message: "user_sign_out"})
}
