package database

import (
	"encoding/json"
	"log"
	"time"

	//	"strconv"

	"drops-backend/models"
	"drops-backend/utils"

	_ "github.com/go-sql-driver/mysql"
)

/*var UserQuery = "SELECT u.id, u.uuid, c.email, u.first_name, u.last_name, u.updated, u.created, CONCAT('[', " +
"GROUP_CONCAT(JSON_OBJECT('uuid', a.uuid, 'role_uuid', r.uuid, 'role_name', r.name, 'service_name', r.service_name, " +
"'model_uuid', m.uuid, 'model_name', m.name, 'created', a.created)), ']') " +
"FROM pool_user AS u " +
"LEFT JOIN credentials AS c ON u.id = c.pool_user_id " +
"LEFT JOIN access_user AS a ON u.id = a.pool_user_id " +
"LEFT JOIN model AS m ON a.id = m.access_user_id " +
"LEFT JOIN role AS r ON a.role_id = r.id "*/

var UserQuery = "SELECT u.id, u.uuid, u.email, u.confirmed, u.updated, u.created, " +
	"p.uuid, p.first_name, p.last_name, ifnull(p.mobile, ''), ifnull(p.birthdate, 0), ifnull(p.gender, 'none'), p.updated, p.created, " +
	"CONCAT('[', GROUP_CONCAT(JSON_OBJECT(" +
	"'uuid', ad.uuid, " +
	"'street', ad.street, " +
	"'additional', ad.additional, " +
	"'zip', ad.zip, " +
	"'country', ad.country, " +
	"'google_id', ad.google_id, " +
	"'updated', ad.updated, " +
	"'created', ad.created)), " +
	"']'), " +
	"CONCAT('[', GROUP_CONCAT(JSON_OBJECT(" +
	"'uuid', a.uuid, " +
	"'access_name', a.name, " +
	"'service_name', m.service_name, " +
	"'model_uuid', m.uuid, " +
	"'model_name', m.name, " +
	"'model_type', m.type, " +
	"'created', m.created)), " +
	"']') " +
	"FROM vca_user AS u " +
	"LEFT JOIN profile AS p ON p.vca_user_id = u.id " +
	"LEFT JOIN address AS ad ON ad.profile_id = p.id " +
	"LEFT JOIN access_user AS a ON u.id = a.vca_user_id " +
	"LEFT JOIN model AS m ON a.model_id = m.id "

/**
 * GET /users
 */
func UserList(page *models.Page, sort string, filter *models.UserFilter) (users []models.User, err error) {
	// define the query
	query := UserQuery +
		"GROUP BY u.id, p.uuid " +
		"ORDER BY " + sort +
		"LIMIT ?, ?"

	// execute query
	rows, err := utils.DB.Query(query, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	//initial dummy varibles
	var accessByte []byte
	var addressByte []byte
	var id int
	// create User and AccessUser
	user := new(models.User)
	profile := new(models.Profile)
	access := new([]models.Access)
	address := new([]models.Address) // convert each row to User
	for rows.Next() {
		//scan row and fill user
		err = rows.Scan(&id, &user.Uuid, &user.Email, &user.Confirmed, &user.Updated, &user.Created,
			&profile.Uuid,
			&profile.FirstName,
			&profile.LastName,
			&profile.Mobile,
			&profile.Birthdate,
			&profile.Gender,
			&profile.Created,
			&profile.Updated,
			&addressByte,
			&accessByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// create json from []byte
		err = json.Unmarshal(accessByte, &access)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		if (*access)[0].Uuid != "" {

			user.Access = *access
		}
		// create json from []byte
		err = json.Unmarshal(addressByte, &address)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		if (*address)[0].Uuid != "" {
			user.Profile.Addresses = *address
		}
		user.Profile = *profile

		// append to list of user
		users = append(users, *user)
	}
	return users, err
}

/**
 * GET /users/:id
 */
func UserById(search string) (users []models.User, err error) {
	// execute the query
	userQuery := UserQuery +
		"WHERE u.uuid = ? " +
		"GROUP BY u.id, p.uuid " +
		"LIMIT 1"
	rows, err := utils.DB.Query(userQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	//initial dummy varibles
	var accessByte []byte
	var addressByte []byte
	var id int
	// create User and AccessUser
	user := new(models.User)
	profile := new(models.Profile)
	access := new([]models.Access)
	address := new([]models.Address) // convert each row to User

	// convert each row to User
	for rows.Next() {
		//scan row and fill user
		err = rows.Scan(&id, &user.Uuid, &user.Email, &user.Confirmed, &user.Updated, &user.Created,
			&profile.Uuid,
			&profile.FirstName,
			&profile.LastName,
			&profile.Mobile,
			&profile.Birthdate,
			&profile.Gender,
			&profile.Created,
			&profile.Updated,
			&addressByte,
			&accessByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// create json from []byte
		err = json.Unmarshal(accessByte, &access)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		if (*access)[0].Uuid != "" {

			user.Access = *access
		}
		// create json from []byte
		err = json.Unmarshal(addressByte, &address)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		if (*address)[0].Uuid != "" {
			user.Profile.Addresses = *address
		}
		user.Profile = *profile
		users = append(users, *user)
	}

	if id == 0 {
		err = utils.ErrorNotFound
		return nil, err
	}
	return users, err
}

/**
 * PUT /users
 */
func UserUpdate(user *models.User) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//update user user
	_, err = tx.Exec("UPDATE vca_user SET updated = ? WHERE uuid = ?", time.Now().Unix(), user.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * DELETE /users
 */
func UserDelete(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	rows, err := tx.Query("SELECT id FROM vca_user WHERE uuid = ?", deleteBody.Uuid)
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
	_, err = tx.Exec("DELETE FROM vca_user WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
