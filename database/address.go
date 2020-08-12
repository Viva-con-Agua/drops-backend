package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

/**
 * select list of Crew
 */
func GetAddressDefaultList(page *models.Page, sort string, filter *models.FilterAddress) (Address []models.AddressExtended, err error) {

	// execute the query
	addressQuery := "SELECT a.uuid, pa.primary, a.street, a.additional, a.zip, a.city, a.country, a.google_id, p.uuid, p.email, CONCAT(p.firstname, ' ', p.lastname) AS fullname, p.mobile, a.Updated, a.Created " +
		"FROM address AS a " +
		"LEFT JOIN profile_has_address AS pa ON a.id = pa.address_id " +
		"LEFT JOIN profile AS p ON pa.profile_id = p.id " +
		"WHERE a.street LIKE ? " +
		"GROUP BY a.id " +
		sort + " " +
		"LIMIT ?, ?"
	rows, err := utils.DB.Query(addressQuery, filter.Street, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	// convert each row
	for rows.Next() {

		//create Profile and corresponding models
		address := new(models.AddressExtended)
		profile := new(models.AddressProfile)

		//scan row and fill Profile
		err = rows.Scan(&address.Uuid, &address.Primary, &address.Street, &address.Additional, &address.Zip, &address.City, &address.Country, &address.GoogleId, &profile.Uuid, &profile.Email, &profile.Fullname, &profile.Mobile, &address.Updated, &address.Created)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		address.Profile = *profile

		Address = append(Address, *address)

	}
	return Address, err
}

/**
 * Create Address
 */
func CreateAddress(Address *models.AddressCreate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing profile
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", Address.ProfileId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}

	// Insert Address
	uuid := uuid.New()
	res, err := tx.Exec("INSERT INTO address (uuid, street, additional, zip, city, country, google_id, updated, created) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		uuid, Address.Street, Address.Additional, Address.Zip, Address.City, Address.Country, Address.GoogleId, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	addressId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// TODO CHECK FOR ANOTHER PRIMARY ADDRESS AND OVERWRITE IF NECCESSARY

	// Insert relation address <> profile
	res, err = tx.Exec("INSERT INTO profile_has_address (profile_id, address_id, `primary`) VALUES "+
		"(?, ?, ?)",
		id, addressId, Address.Primary)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

/**
 * select Profile by uuid
 */
func GetAddress(search string) (Address []models.AddressExtended, err error) {

	// execute the query
	addressQuery := "SELECT a.uuid, a.street, a.additional, a.zip, a.city, a.country, a.google_id, p.uuid, p.email, CONCAT(p.firstname, ' ', p.lastname) AS fullname, p.mobile, a.Updated, a.Created " +
		"FROM address AS a " +
		"LEFT JOIN profile_has_address AS pa ON a.id = pa.address_id " +
		"LEFT JOIN profile AS p ON pa.profile_id = p.id " +
		"WHERE a.uuid = ?"

	rows, err := utils.DB.Query(addressQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	// convert each row
	for rows.Next() {

		//create Profile and corresponding models
		address := new(models.AddressExtended)
		profile := new(models.AddressProfile)

		//scan row and fill Profile
		err = rows.Scan(&address.Uuid, &address.Street, &address.Additional, &address.Zip, &address.City, &address.Country, &address.GoogleId, &profile.Uuid, &profile.Email, &profile.Fullname, &profile.Mobile, &address.Updated, &address.Created)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		address.Profile = *profile

		Address = append(Address, *address)
	}
	return Address, err
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
	// TODO CHECK FOR ANOTHER PRIMARY ADDRESS AND OVERWRITE IF NECCESSARY

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

func DeleteAddress(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// select id
	rows, err := tx.Query("SELECT id FROM address WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}

	if id == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Delete address relation
	_, err = tx.Exec("DELETE FROM profile_has_address WHERE address_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete address itself
	// TODO DELETE PROFILE AND CORRESPONDING RELATIONS
	_, err = tx.Exec("DELETE FROM address WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}
