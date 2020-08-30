package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/Viva-con-Agua/echo-pool/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
 * TODO more than on service assign
 */
func SignUp(signup_data *models.SignUpData) (user_uuid *string, access_token *string, err error) {
	// select service
	rows, err := utils.DB.Query("SELECT s.id FROM service As s WHERE s.name='drops-backend'")
	if err != nil {
		log.Print(err, " ### database.SignUp Step_1")
		return nil, nil, err
	}
	// select model_id from rows
	var service_id int
	for rows.Next() {
		err = rows.Scan(&service_id)
		if err != nil {
			log.Print(err, " ### database.SignUp Step_2")
			return nil, nil, err
		}
	}

	// begin database query and handle error
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print(err, " ### database.SignUp Step_3")
		return nil, nil, err
	}
	//insert user
	u_uuid := uuid.New().String()
	res, err := tx.Exec(
		"INSERT INTO drops_user (uuid, email, confirmed, updated, created) VALUES(?, ?, ?, ?, ?)",
		u_uuid,
		signup_data.SignUpUser.Email,
		0,
		time.Now().Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, nil, utils.ErrorConflict
		}
		log.Print(err, " ### database.SignUp Step_4")
		return nil, nil, err
	}
	// get user id via LastInsertId
	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err, " ### database.SignUp Step_5")
		return nil, nil, err
	}

	// insert credentials
	password, err := bcrypt.GenerateFromPassword([]byte(signup_data.SignUpUser.Password), 10)
	res, err = tx.Exec(
		"INSERT INTO password_info (password, hasher, drops_user_id) VALUES(?, ?, ?)",
		password,
		"bcrypt",
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_6")
		return nil, nil, err
	}
	//insert profile
	// Create uuid
	profile_uuid := uuid.New()
	res, err = tx.Exec(
		"INSERT INTO profile (uuid, first_name, last_name, updated, created, drops_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		profile_uuid,
		signup_data.SignUpUser.FirstName,
		signup_data.SignUpUser.LastName,
		time.Now().Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_7")
		return nil, nil, err
	}

	// begin database query and handle error
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, nil, err
	}
	// Create uuid
	default_access_uuid := uuid.New()
	_, err = tx.Exec("INSERT INTO access (uuid, name, created, service_id, drops_user_id ) VALUES(?, ?, ?, ?, ?)",
		default_access_uuid,
		"joined",
		time.Now().Unix(),
		service_id,
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_8")
		return nil, nil, err
	}
	//insert access_token
	token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_10")
		return nil, nil, err
	}
	res, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, redirect_url, expired, created,drops_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		token,
		"signup",
		signup_data.RedirectUrl,
		time.Now().Add(time.Hour*24).Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_11")
		return nil, nil, err
	}

	// insert profile
	return &u_uuid, &token, tx.Commit()
}

func ConfirmSignUp(t string) (user_uuid *string, err error) {
	query := "SELECT uuid FROM drops_user " +
		"JOIN access_token ON access_token.drops_user_id = drops_user.id " +
		"WHERE access_token.token = ? && access_token.expired > ?"
	rows, err := utils.DB.Query(query, t, time.Now().Unix())
	var u_uuid string
	for rows.Next() {
		err = rows.Scan(&u_uuid)
		if err != nil {
			log.Print(err, " ### database.ConfirmSignUp Step_1")
			return nil, err
		}
	}

	query = "UPDATE drops_user " +
		"JOIN access_token ON access_token.drops_user_id = drops_user.id " +
		"SET drops_user.confirmed = 1, updated = ?  " +
		"WHERE access_token.token = ? && access_token.expired > ?"
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print(err, " ### database.ConfirmSignUp Step_2")
		return nil, err
	}
	//update user user
	res, err := tx.Exec(query, time.Now().Unix(), t, time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print(": ", err, " ### database.ConfirmSignUp Step_3")
		return nil, err
	}
	r, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.ConfirmSignUp Step_4 ")
		return nil, err
	}
	if r == 0 {
		err = utils.ErrorNotFound
		return nil, err
	}
	_, err = tx.Exec("DELETE FROM access_token WHERE access_token.token = ?", t)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.ConfirmSignUp Step_5 ")
		return nil, err
	}
	return &u_uuid, tx.Commit()
}

func SignUpToken(n *models.NewToken) (*string, error) {

	//insert access_token
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}
	query := "SELECT drops_user.id FROM drops_user " +
		"JOIN password_info ON drops_user.id = password_info.drops_user_id " +
		"WHERE drops_user.email = ? && drops_user.confirmed = 0 " +
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
		"INSERT INTO access_token (token, t_case, expired, created, drops_user_id) VALUES(?, ?, ?, ?, ?)",
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

func GetSessionUser(user_uuid *string) (user *auth.User, err error) {
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
		"'access_uuid', a.uuid, " +
		"'access_name', a.name, " +
		"'service_name', s.name, " +
		"'model_uuid', m.uuid, " +
		"'model_name', m.name, " +
		"'model_type', m.type, " +
		"'created', a.created)), " +
		"']') " +
		"FROM drops_user AS u " +
		"LEFT JOIN profile AS p ON p.drops_user_id = u.id " +
		"LEFT JOIN address AS ad ON ad.profile_id = p.id " +
		"LEFT JOIN access AS a ON u.id = a.drops_user_id " +
		"LEFT JOIN service AS s ON a.service_id = s.id " +
		"LEFT JOIN access_has_model AS ahs ON ahs.access_id = a.id " +
		"LEFT JOIN model AS m ON m.id = ahs.model_id " +
		"WHERE u.uuid = ? " +
		"GROUP BY u.id, u.email, p.uuid " +
		"LIMIT 1"
	rows, err := utils.DB.Query(query, user_uuid)
	if err != nil {
		log.Print(err, " ### database.GetSessionUser Step_1")
		return nil, err
	}
	//initial dummy varibles
	var accessByte []byte
	var addressByte []byte
	var id int
	// create User and AccessUser
	user = new(auth.User)
	profile := new(auth.Profile)
	access := new(models.AccessDBList)
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
		//if (*access)[0].AccessUuid != "" {

		user.Access = *access.AccessList()
		//}
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

func SignIn(c *models.SignInData) (user *auth.User, err error) {

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
		"'created', a.created)), " +
		"']'), " +
		"pi.password " +
		"FROM drops_user AS u " +
		"LEFT JOIN profile AS p ON p.drops_user_id = u.id " +
		"LEFT JOIN address AS ad ON ad.profile_id = p.id " +
		"LEFT JOIN password_info AS pi ON pi.drops_user_id = u.id " +
		"LEFT JOIN access AS a ON u.id = a.drops_user_id " +
		"LEFT JOIN service AS s ON a.service_id = s.id " +
		"LEFT JOIN access_has_model AS ahs ON ahs.access_id = a.id " +
		"LEFT JOIN model AS m ON m.id = ahs.model_id " +
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
	access := new(models.AccessDBList)
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
		if (*access)[0].AccessUuid != "" {

			user.Access = *access.AccessList()
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
