package database

import (
    "../models"
    "../utils"
    _ "github.com/go-sql-driver/mysql"
    "github.com/google/uuid"
    "log"
    "time"
)

/**
 * select addresses for an given profile_id
 */
func GetAddressesByProfileId(profile_id int) (addresses []models.Address, err error) {
    // Execute the Query
    // TODO ADD PRIMARY INFORMATIONS
    rows, err := utils.DB.Query("SELECT address.uuid, address.street, address.additional, address.zip, address.city, address.country, address.google_id, address.updated, address.created FROM address LEFT JOIN profile_has_address ON profile_has_address.address_id = address.id WHERE profile_has_address.profile_id = ?", profile_id)
    if err != nil {
        log.Print("Database Error", err)
        return nil, err
    }
    // define variable for each column
    // convert each row
    for rows.Next() {
        address := new(models.Address)
        //scan row
        err = rows.Scan(&address.Uuid, &address.Street, &address.Additional, &address.Zip, &address.City, &address.Country, &address.GoogleId, &address.Updated, &address.Created)
        if err != nil {
            log.Print("Database Error: ", err)
            return nil, err
        }
        // fill []models.Address
        addresses = append(addresses, *address)
    }
    return addresses, err
}

/**
 * Create Address
 */
func CreateAddress(address *models.AddressCreate) (err error) {
    // sgl begin
    tx, err := utils.DB.Begin()
    if err != nil {
        log.Print("Database Error: ", err)
        return err
    }

    // Insert Profile
    uuid := uuid.New()
    _, err = tx.Exec("INSERT INTO address (uuid, street, additional, zip, city, country, google_id, updated, created) VALUES "+
        "(?, ?, ?, ?, ?, ?, ?, ?, ?)",
        uuid, address.Street, address.Additional, address.Zip, address.City, address.Country, address.GoogleId, time.Now().Unix(), time.Now().Unix())
    if err != nil {
        tx.Rollback()
        log.Print("Database Error: ", err)
        return err
    }
    return tx.Commit()

}

/**
 * update Profile
 */
func UpdateAddress(Address *models.AddressUpdate) (err error) {
    // sgl begin
    tx, err := utils.DB.Begin()
    if err != nil {
        log.Print("Database Error: ", err)
        return err
    }
    //slect id
    rows, err := tx.Query("SELECT id FROM address WHERE uuid = ?", Address.Uuid)
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

    //update profile
    _, err = tx.Exec("UPDATE address SET street = ?, additional = ?, zip = ?, city = ?, country = ?, google_id = ?, updated = ? WHERE id = ?", Address.Street, Address.Additional, Address.Zip, Address.City, Address.Country, Address.GoogleId, time.Now().Unix(), id)
    if err != nil {
        tx.Rollback()
        log.Print("Database Error: ", err)
        return err
    }

    // TODO UPDATE ADDRESSES

    return tx.Commit()

}
