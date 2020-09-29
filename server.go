package main

import (
	"drops-backend/controllers"
	"drops-backend/nats"
	"drops-backend/utils"
	"log"
	"os"
	"strings"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
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
	godotenv.Load()
	log.Print(strings.Split(os.Getenv("ALLOW_ORIGINS"), ","))
	utils.ConnectDatabase()
	store := api.RedisSession()
	nats.Connect()
	nats.SubscribeAddModel()
	//create echo server
	e := echo.New()
	m := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(os.Getenv("ALLOW_ORIGINS"), ","),
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	})
	e.Use(m)
	e.Use(store)
	e.Validator = &CustomValidator{validator: validator.New()}

	apiV1 := e.Group("/v1")

	// "/v1/auth"
	a := apiV1.Group("/auth")
	a.POST("/signup", controllers.SignUp)
	a.GET("/signup/confirm/:token", controllers.ConfirmSignUp)
	a.POST("/signin", controllers.SignIn)
	a.POST("/signup/token", controllers.SignUpToken)
	a.GET("/current", controllers.Current)
	a.GET("/signout", controllers.SignOut)

	// "/v1/users"
	users := apiV1.Group("/users")
	users.Use(api.SessionAuth)
	users.POST("/password", controllers.PasswordResetToken)
	users.PUT("/password", controllers.PasswordReset)
	users.DELETE("", controllers.UserDelete)

	// "/v1/access"
	apiV1.POST("/access", controllers.AccessInsert)
	apiV1.DELETE("/access", controllers.AccessDelete)

	// "v1/models"
	apiV1.POST("/models", controllers.ModelCreate)
	apiV1.DELETE("/models", controllers.ModelDelete)

	apiAdmin := e.Group("/admin")

	apiAdmin.DELETE("/users", controllers.UserDelete)

	// "/v1/access"
	apiAdmin.POST("/access", controllers.AccessInsert)
	apiAdmin.DELETE("/access", controllers.AccessDelete)

	// "v1/models"
	apiAdmin.POST("/models", controllers.ModelCreate)
	apiAdmin.DELETE("/models", controllers.ModelDelete)

	//internal routes for microservices
	//intern := e.Group("/intern")
	if port, ok := os.LookupEnv("REPO_PORT"); ok {
		e.Logger.Fatal(e.Start(":" + port))
	} else {
		e.Logger.Fatal(e.Start(":1323"))
	}
}
