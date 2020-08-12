package main

import (
	"drops-backend/controllers"
	"drops-backend/utils"
	"strconv"

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
	store := pool.RedisSession(utils.Config.REDIS.Host + ":" + strconv.Itoa(utils.Config.REDIS.Port))

	//create echo server
	e := echo.New()
	e.Use(store)
	e.Validator = &CustomValidator{validator: validator.New()}

	// TODO: Listen for user creation on nats

	apiV1 := e.Group("/v1/drops")
	// TODO REENABLE AUTHENTICATION
	//apiV1.Use(pool.SessionAuth)

	// ROUTES FOR PROFILES
	apiV1.GET("/profiles", controllers.GetProfileDefaultList)
	apiV1.GET("/profiles/ids", controllers.GetProfileIdList)
	apiV1.GET("/profiles/min", controllers.GetProfileMinList)
	apiV1.GET("/profiles/extended", controllers.GetProfileExtendedList)

	apiV1.POST("/profile", controllers.CreateProfile)
	apiV1.GET("/profile/:id", controllers.ReadProfile)
	apiV1.PUT("/profile", controllers.UpdateProfile)
	apiV1.DELETE("/profile", controllers.DeleteProfile)

	// ROUTE FOR NEWSLETTER SELECTION
	apiV1.POST("/newsletter", controllers.SetNewsletter)

	// ROUTES FOR AVATAR
	// TODO TEST UPLOADS WITH FILE
	apiV1.PUT("/avatar", controllers.SetAvatar)
	apiV1.DELETE("/avatar", controllers.DeleteAvatar)

	apiV1.GET("/addresses", controllers.GetAddressDefaultList)
	apiV1.POST("/address", controllers.CreateAddress)
	apiV1.GET("/address/:id", controllers.ReadAddress)
	apiV1.PUT("/address", controllers.UpdateAddress)
	apiV1.DELETE("/address", controllers.DeleteAddress)

	// ROUTES FOR CREWS
	apiV1.GET("/crews", controllers.GetCrewList)

	apiV1.POST("/crew", controllers.CreateCrew)
	apiV1.GET("/crew/:id", controllers.ReadCrew)
	apiV1.PUT("/crew", controllers.UpdateCrew)
	apiV1.DELETE("/crew", controllers.DeleteCrew)

	// ROUTES FOR CREW SELECTION
	apiV1.POST("/crew/assign", controllers.AssignCrew)
	apiV1.DELETE("/crew/remove", controllers.RemoveCrew)

	// ROUTES FOR ROLES
	apiV1.GET("/roles", controllers.GetRolesDefaultList)

	apiV1.POST("/role", controllers.CreateRole)
	apiV1.GET("/role/:id", controllers.ReadRole)
	apiV1.PUT("/role", controllers.UpdateRole)
	apiV1.DELETE("/role", controllers.DeleteRole)

	// ROUTES FOR ASP ASSIGNMENT
	apiV1.POST("/role/assign", controllers.AssignRole)
	apiV1.DELETE("/role/remove", controllers.RemoveRole)

	apiV1.POST("/active/:state", controllers.ActiveStateChange)

	apiV1.POST("/nvm/:state", controllers.NVMStateChange)

	// TODO: ADD NVM_STATE
	// TODO V2: Permission assignment

	// TODO: ADD ROUTES FOR ASP ASSIGNMENT
	// TODO: ADD ROUTES FOR AVATARS

	e.Logger.Fatal(e.Start(":1323"))
}
