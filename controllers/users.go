package controllers

import (
	"log"
	"net/http"

	"drops-backend/database"
	"drops-backend/models"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

/**
 * "GET /users"
 * response list of models.User
 */
func UserList(c echo.Context) (err error) {
	query := new(models.UserQuery)
	if err = c.Bind(query); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	page := query.Page()
	sort := query.OrderBy()
	filter := query.Filter()
	response, err := database.UserList(page, sort, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func UserListInternal(c echo.Context) (err error) {
	filter := new(api.UserRequest)
	if err = c.Bind(filter); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	response, err := database.UserListInternal(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
* "GET /users/:uuid"
* response models.User
 */
func UserById(c echo.Context) (err error) {
	uuid := c.Param("uuid")
	response, err := database.UserById(uuid)
	if err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

/**
 * "GET /users"
 * response uuid of updated models.User
 */
func UserUpdate(c echo.Context) (err error) {
	// create body as models.User
	body := new(api.User)
	// save data to body
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// validate body
	if err = c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// update body into database
	if err = database.UserUpdate(body); err != nil {
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	//TODO nats.update.user
	return c.JSON(http.StatusOK, resp.Updated(body.Uuid))
}

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
		if err == utils.ErrorNotFound {
			return c.JSON(http.StatusNoContent, resp.NoContent(body.Uuid))
		}
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError())
	}
	// response created
	return c.JSON(http.StatusOK, resp.Deleted(body.Uuid))
}

func ProfileDelete(u_uuid string) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print(err, " ### database.ProfileDelete")
		return err
	}
	rows, err := tx.Query("SELECT id FROM drops_user WHERE uuid = ?", u_uuid)
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFound
	if id == 0 {
		err = utils.ErrorNotFound
		return err
	}
	//update user user
	_, err = tx.Exec("DELETE FROM drops_user WHERE uuid = ?", u_uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
