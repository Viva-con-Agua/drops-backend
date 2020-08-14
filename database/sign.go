package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/Viva-con-Agua/echo-pool/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *models.SignUpData) (*string, error) {

	// Create uuid
	Uuid, err := uuid.NewRandom()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}

	// begin database query and handle error
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}
	//insert user
	res, err := tx.Exec(
		"INSERT INTO vca_user (uuid, email, confirmed, updated, created) VALUES(?, ?, ?, ?, ?)",
		Uuid.String(),
		c.Email,
		0,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	// get user id via LastInsertId
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}

	// insert credentials
	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), 10)
	res, err = tx.Exec(
		"INSERT INTO password_info (password, hasher, vca_user_id) VALUES(?, ?, ?)",
		password,
		"bcrypt",
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	//insert profile
	// Create uuid
	Uuid, err = uuid.NewRandom()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}
	res, err = tx.Exec(
		"INSERT INTO profile (uuid, first_name, last_name, updated, created, vca_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		Uuid,
		c.FirstName,
		c.LastName,
		time.Now().Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}

	//insert access_token
	access_token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	res, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, expired, created, vca_user_id) VALUES(?, ?, ?, ?, ?)",
		access_token,
		"signup",
		time.Now().Add(time.Hour*24).Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}

	// insert profile
	return &access_token, tx.Commit()
}

func ConfirmSignUp(t string) (err error) {
	query := "UPDATE vca_user " +
		"JOIN access_token ON access_token.vca_user_id = vca_user.id " +
		"SET vca_user.confirmed = 1, updated = ?  " +
		"WHERE access_token.token = ? && access_token.expired > ?"
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//update user user
	res, err := tx.Exec(query, time.Now().Unix(), t, time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Print("database.ConfirmSignUp: ", err)
		return err
	}
	if r == 0 {
		err = utils.ErrorNotFound
		return err
	}
	_, err = tx.Exec("DELETE FROM access_token WHERE access_token.token = ?", t)
	if err != nil {
		tx.Rollback()
		log.Print("database.ConfirmSignUp: ", err)
		return err
	}
	return tx.Commit()
}

func SignUpToken(n *models.NewToken) (*string, error) {

	//insert access_token
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}
	query := "SELECT vca_user.id FROM vca_user " +
		"JOIN password_info ON vca_user.id = password_info.vca_user_id " +
		"WHERE vca_user.email = ? && vca_user.confirmed = 0 " +
		"LIMIT 1"
	rows, err := utils.DB.Query(query, n.Email)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	//initial dummy varibles
	var id int
	// convert each row to User
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
	}
	if id == 0 {
		return nil, utils.ErrorNotFound
	}
	access_token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	_, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, expired, created, vca_user_id) VALUES(?, ?, ?, ?, ?)",
		access_token,
		"signup",
		time.Now().Add(time.Hour*24).Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	return &access_token, tx.Commit()

}

func SignIn(c *models.Credentials) (user *auth.User, err error) {

	query := "SELECT u.id, u.uuid, u.email, u.confirmed, u.updated, u.created, " +
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
		"']'), " +
		"pi.password " +
		"FROM vca_user AS u " +
		"LEFT JOIN profile AS p ON p.vca_user_id = u.id " +
		"LEFT JOIN address AS ad ON ad.profile_id = p.id " +
		"LEFT JOIN password_info AS pi ON pi.vca_user_id = u.id " +
		"LEFT JOIN access_user AS a ON u.id = a.vca_user_id " +
		"LEFT JOIN model AS m ON a.model_id = m.id " +
		"WHERE u.email = ? " +
		"GROUP BY u.id, u.email, p.uuid, pi.password " +
		"LIMIT 1"

	rows, err := utils.DB.Query(query, c.Email)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	//initial dummy varibles
	var accessByte []byte
	var addressByte []byte
	var id int
	var password []byte
	// create User and AccessUser
	user = new(auth.User)
	profile := new(auth.Profile)
	access := new([]auth.Access)
	address := new([]auth.Address)

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
			&accessByte,
			&password)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		err = bcrypt.CompareHashAndPassword(password, []byte(c.Password))
		if err != nil {
			return nil, utils.ErrorPassword
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
	}
	if user.Uuid == "" {
		return nil, err
	}
	return user, err

}

/*func Current(c *models.Credentials) (user *auth.User, err error) {
	query := "SELECT u.id, u.uuid, c.email, u.confirmed, u.updated, u.created, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', a.uuid, 'role_uuid', r.uuid, 'role_name', r.name, 'service_name', r.service_name, " +
		"'model_uuid', m.uuid, 'model_name', m.name, 'created', a.created)), ']'), c.password " +
		"FROM pool_user AS u " +
		"LEFT JOIN credentials AS c ON u.id = c.pool_user_id " +
		"LEFT JOIN access_user AS a ON u.id = a.pool_user_id " +
		"LEFT JOIN model AS m ON a.id = m.access_user_id " +
		"LEFT JOIN role AS r ON a.role_id = r.id " +
		"WHERE c.email = ? " +
		"GROUP BY u.id, c.email " +
		"LIMIT 1"
	rows, err := utils.DB.Query(query, c.Email)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	//initial dummy varibles
	var accessByte []byte
	var id int
	var password []byte
	// create User and AccessUser
	user = new(auth.User)
	access := new([]auth.AccessUser)

	// convert each row to User
	for rows.Next() {

		//scan row and fill user
		err = rows.Scan(&id, &user.Uuid, &user.Email, &user.Confirmed, &user.Updated, &user.Created, &accessByte, &password)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		err = bcrypt.CompareHashAndPassword(password, []byte(c.Password))
		if err != nil {
			return nil, utils.ErrorPassword
		}
		// create json from []byte
		err = json.Unmarshal(accessByte, &access)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		user.Access = *access
	}
	if user.Uuid == "" {
		return nil, err
	}
	return user, err

}*/
