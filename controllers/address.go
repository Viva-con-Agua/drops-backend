package controllers

import (
    "../database"
    "../models"
    "../utils"
    "github.com/Viva-con-Agua/echo-pool/pool"
    "github.com/labstack/echo"
    "net/http"
)

/**
 * Response list of models.AddressList
 */
func GetAddressDefaultList(c echo.Context) (err error) {
    query := new(models.QueryAddress)
    if err = c.Bind(query); err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }
    query.Defaults()
    page := query.Page()
    sort := query.OrderBy()
    filter := query.Filter()
    response, err := database.GetAddressDefaultList(page, sort, filter)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
    }
    return c.JSON(http.StatusOK, response)
}

func CreateAddress(c echo.Context) (err error) {
    // create body as models.ProfileCreate
    body := new(models.AddressCreate)
    // save data to body
    if err = c.Bind(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    // validate body
    if err = c.Validate(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    // update body into database
    if err = database.CreateAddress(body); err != nil {
        if err == utils.ErrorConflict {
            return c.JSON(http.StatusNoContent, pool.Conflict())
        }
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
    }
    // response created
    return c.JSON(http.StatusOK, pool.Created())
}

/**
 * Response list of models.Address
 */
func ReadAddress(c echo.Context) (err error) {
    uuid := c.Param("id")
    response, err := database.GetAddress(uuid)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError)
    }
    if response == nil {
        return c.JSON(http.StatusNoContent, pool.NoContent(uuid))
    }
    return c.JSON(http.StatusOK, response)
}

func UpdateAddress(c echo.Context) (err error) {
    // create body as models.Profile
    body := new(models.AddressUpdate)

    // save data to body
    if err = c.Bind(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    // validate body
    if err = c.Validate(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    // update body into database
    if err = database.UpdateAddress(body); err != nil {
        if err == utils.ErrorNotFound {
            return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
        }
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
    }
    // response created
    return c.JSON(http.StatusOK, pool.Updated(body.Uuid))
}

func DeleteAddress(c echo.Context) (err error) {
    // create body as models.DeleteBody
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
    if err = database.DeleteAddress(body); err != nil {
        if err == utils.ErrorNotFound {
            return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
        }
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
    }
    // response created
    return c.JSON(http.StatusOK, pool.Deleted(body.Uuid))
}
