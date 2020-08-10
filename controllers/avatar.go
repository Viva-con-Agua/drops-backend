package controllers

import (
    "../database"
    "../models"
    "../utils"
    "github.com/Viva-con-Agua/echo-pool/pool"
    "github.com/labstack/echo"
    "net/http"
)

func SetAvatar(c echo.Context) (err error) {
    // create body as models.Profile
    body := new(models.AvatarUpdate)

    // save data to body
    if err = c.Bind(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }
    // validate body
    if err = c.Validate(body); err != nil {
        return c.JSON(http.StatusBadRequest, err)
    }

    file, err := c.FormFile("file")
    if err != nil {
        return err
    }
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    // update body into database
    if err = database.UpdateAvatar(body, src); err != nil {
        if err == utils.ErrorNotFound {
            return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
        }
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
    }
    // response created
    return c.JSON(http.StatusOK, pool.Updated(body.Uuid))
}

func DeleteAvatar(c echo.Context) (err error) {
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
    if err = database.DeleteAvatar(body); err != nil {
        if err == utils.ErrorNotFound {
            return c.JSON(http.StatusNoContent, pool.NoContent(body.Uuid))
        }
        return c.JSON(http.StatusInternalServerError, pool.InternelServerError())
    }
    // response created
    return c.JSON(http.StatusOK, pool.Deleted(body.Uuid))
}
