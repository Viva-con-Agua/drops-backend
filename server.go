package main

import (
	"drops-backend/controllers"
	"drops-backend/nats"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/auth"
	"github.com/Viva-con-Agua/echo-pool/config"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	// intial loading function
	utils.LoadConfig()
	config.LoadConfig()
	utils.ConnectDatabase()
	store := auth.RedisSession()
	nats.Connect()
	//create echo server
	e := echo.New()
	m := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     utils.Config.Alloworigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	})
	e.Use(m)
	e.Use(store)
	e.Validator = &CustomValidator{validator: validator.New()}

	a := e.Group("/auth")
	a.POST("/signup", controllers.SignUp)
	a.GET("/signup/confirm/:token", controllers.ConfirmSignUp)
	a.POST("/signin", controllers.SignIn)
	a.POST("/signup/token", controllers.SignUpToken)
	a.GET("/current", controllers.Current)

	a.GET("/signout", controllers.SignOut)

	apiV1 := e.Group("/v1")
	apiV1.Use(auth.SessionAuth)

	// "/v1/users"
	apiV1.GET("/users/:uuid", controllers.UserById)
	apiV1.GET("/users", controllers.UserList)
	apiV1.PUT("/users", controllers.UserUpdate)
	apiV1.DELETE("/users", controllers.UserDelete)

	// "/v1/access"
	apiV1.POST("/access", controllers.AccessInsert)
	apiV1.DELETE("/access", controllers.AccessDelete)

	// "v1/models"
	apiV1.POST("/models", controllers.ModelInsert)
	apiV1.DELETE("/models", controllers.ModelDelete)

	e.Logger.Fatal(e.Start(":1323"))
}
