package main

import (
    "./controllers"
    "./utils"
    "github.com/Viva-con-Agua/echo-pool/pool"
    "github.com/go-playground/validator"
    "github.com/labstack/echo"
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
    utils.ConnectDatabase()
    store := pool.RedisSession("172.2.150.2:6379")

    //create echo server
    e := echo.New()
    e.Use(store)
    e.Validator = &CustomValidator{validator: validator.New()}

    // TODO: Listen for user creation on nats

    apiV1 := e.Group("/v1/")
    // TODO REENABLE AUTHENTICATION
    //apiV1.Use(pool.SessionAuth)

    apiV1.GET("/profiles", controllers.GetProfileDefaultList)
    apiV1.GET("/profiles/ids", controllers.GetProfileIdList)
    apiV1.GET("/profiles/min", controllers.GetProfileMinList)
    apiV1.GET("/profiles/extended", controllers.GetProfileExtendedList)

    apiV1.POST("/profile", controllers.CreateProfile)
    apiV1.GET("/profile/:id", controllers.GetProfile)
    apiV1.PUT("/profile", controllers.UpdateProfile)
    apiV1.DELETE("/profile", controllers.DeleteProfile)

    //apiV1.GET("/addresses/:id", controllers.GetAdressList)
    //apiV1.GET("/address/:id", controllers.GetAdress)
    //apiV1.POST("/address", controllers.CreateAddress)
    //apiV1.PUT("/address", controllers.UpdateAddress)
    //apiV1.DELETE("/address", controllers.DeleteAddress)

    // TODO: UPDATE ROUTES FOR ENTITIES
    apiV1.GET("/crew", controllers.GetCrewList)
    //apiV1.GET("/crew/:id", controllers.GetCrew)
    apiV1.PUT("/crew", controllers.UpdateCrew)
    apiV1.DELETE("/crew", controllers.DeleteCrew)
    //apiV1.POST("/crew", controllers.CreateCrew)

    // TODO: ADD ROUTES FOR ASP ASSIGNMENT
    // TODO: ADD ROUTES FOR AVATARS

    e.Logger.Fatal(e.Start(":1323"))
}
