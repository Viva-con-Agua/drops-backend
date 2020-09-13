package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
* TODO more than on service assign
 */
func SignUp(s *models.SignUp) (user_uuid *string, access_token *string, err_api *api.ApiError) {
	// begin database query and handle error
	tx, err := utils.DB.Begin()
	if err != nil {
		return nil, nil, api.GetError(err)
	}
	created := time.Now().Unix()
	//insert user
	u_uuid := uuid.New().String()
	res, err := tx.Exec(
		"INSERT INTO drops_user (uuid, email, confirmed, privacy_policy, updated, created) VALUES(?, ?, ?, ?, ?, ?)",
		u_uuid,
		s.SignUser.Email,
		0,
		s.SignUser.PrivacyPolicy,
		created,
		created,
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	// get user id via LastInsertId
	u_id, err := res.LastInsertId()
	if err != nil {
		return nil, nil, api.GetError(err)
	}
	res, err = tx.Exec(
		"INSERT INTO profile (uuid, first_name, last_name, updated, created, drops_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		uuid.New(),
		s.SignUser.FirstName,
		s.SignUser.LastName,
		created,
		created,
		u_id,
	)
	if err != nil {
		tx.Rollback()
		log.Print(err, " ### database.SignUp Step_7")
		return nil, nil, api.GetError(err)
	}
	// insert credentials
	password, err := bcrypt.GenerateFromPassword([]byte(s.SignUser.Password), 10)
	res, err = tx.Exec(
		"INSERT INTO password_info (password, hasher, updated, created, drops_user_id) VALUES(?, ?, ?, ?, ?)",
		password,
		"bcrypt",
		created,
		created,
		u_id,
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	// Access
	query := "INSERT INTO access (uuid, service, updated, created, drops_user_id) " +
		"VALUES(?, 'drops-backend', ?, ?, ?)"
	default_access_uuid := uuid.New()
	_, err = tx.Exec(
		query,
		default_access_uuid,
		created,
		created,
		u_id,
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	a_id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	query = "INSERT INTO `access_right` ( `right`, created, access_id) " +
		"VALUES('joined', ?, ?)"
	_, err = tx.Exec(query, created, a_id)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	//insert access_token
	token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	//Token
	res, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, redirect_url, expired, created,drops_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		token,
		"signup",
		s.RedirectUrl,
		time.Now().Add(time.Hour*24).Unix(),
		created,
		u_id,
	)
	if err != nil {
		tx.Rollback()
		return nil, nil, api.GetError(err)
	}
	// insert profile
	return &u_uuid, &token, api.GetError(tx.Commit())
}

func ConfirmSignUp(t string) (u_uuid *string, api_err *api.ApiError) {
	created := time.Now().Unix()
	var u_id int64
	query := "SELECT u.id, u.uuid FROM drops_user AS u " +
		"JOIN access_token AS t ON t.drops_user_id = u.id " +
		"WHERE t.token = ? && t.expired > ?"
	err := utils.DB.QueryRow(query, t, created).Scan(&u_id, &u_uuid)
	if err != nil {
		return nil, api.GetError(err)
	}
	query = "UPDATE drops_user " +
		"JOIN access_token ON access_token.drops_user_id = drops_user.id " +
		"SET drops_user.confirmed = 1, updated = ?  " +
		"WHERE access_token.token = ? && access_token.expired > ?"
	tx, err := utils.DB.Begin()
	if err != nil {
		return nil, api.GetError(err)
	}
	//update user user
	res, err := tx.Exec(query, created, t, created)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	r, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	if r == 0 {
		err = utils.ErrorNotFound
		return nil, api.GetError(err)
	}

	query = "INSERT INTO `access_right` ( `right`, created, access_id) " +
		"VALUES('created', ?, (SELECT a.id FROM access AS a WHERE a.service = 'drops-backend' AND a.model_id IS NULL AND a.drops_user_id = ?))"
	_, err = tx.Exec(query, created, u_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	_, err = tx.Exec("DELETE FROM access_token WHERE access_token.token = ?", t)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	return u_uuid, api.GetError(tx.Commit())
}

//TODO ResetPassword

func SignUpToken(n *models.NewToken) (t *string, ap_err *api.ApiError) {
	query := "SELECT drops_user.id FROM drops_user " +
		"WHERE drops_user.email = ? && drops_user.confirmed = 0 " +
		"LIMIT 1"
	var u_id int64
	err := utils.DB.QueryRow(query, n.Email).Scan(&u_id)
	if err != nil {
		return nil, api.GetError(err)
	}
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, api.GetError(err)
	}
	access_token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	_, err = tx.Exec(
		"UPDATE access_token SET token = ?, expired = ?, created = ? WHERE t_case = 'signup' AND drops_user_id = ?",
		access_token,
		time.Now().Add(time.Hour*24).Unix(),
		time.Now().Unix(),
		u_id,
	)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	return &access_token, api.GetError(tx.Commit())
}

func GetSessionUser(user_uuid *string) (user *api.UserSession, err_api *api.ApiError) {
	query := "SELECT u.uuid, u.email, u.confirmed, u.privacy_policy, u.updated, u.created, " +
		"p.uuid, u.uuid, p.first_name, p.last_name, CONCAT(p.first_name, ' ', p.last_name), " +
		"ifnull(p.display_name, 'none'), ifnull(p.gender, 'none'),p.updated, p.created, " +
		"ifnull(av.type, ''), ifnull(av.url,''), ifnull(av.updated, 0), ifnull(av.created, 0), " +
		"CONCAT('[', GROUP_CONCAT(JSON_OBJECT(" +
		"'service', a.service, " +
		"'model', 'default', " +
		"'right', ar.right )), " +
		"']') " +
		"FROM drops_user AS u " +
		"LEFT JOIN profile AS p ON p.drops_user_id = u.id " +
		"LEFT JOIN avatar AS av ON av.profile_id = av.id " +
		"LEFT JOIN access AS a ON u.id = a.drops_user_id " +
		"LEFT JOIN access_right AS ar ON ar.access_id = a.id " +
		"LEFT JOIN model AS m ON a.model_id = m.id " +
		"WHERE u.uuid = ? " +
		"GROUP BY u.id, p.uuid, av.id " +
		"LIMIT 1"
	var accessByte []byte
	user = new(api.UserSession)
	profile := new(models.Profile)
	err := utils.DB.QueryRow(query, user_uuid).Scan(
		&user.Uuid,
		&user.Email,
		&user.Confirmed,
		&user.PrivacyPolicy,
		&user.Updated,
		&user.Created,
		&profile.Uuid,
		&profile.UserUuid,
		&profile.FirstName,
		&profile.LastName,
		&profile.FullName,
		&profile.DisplayName,
		&profile.Gender,
		&profile.Created,
		&profile.Updated,
		&profile.Avatar.Type,
		&profile.Avatar.Url,
		&profile.Avatar.Updated,
		&profile.Avatar.Created,
		&accessByte)
	if err != nil {
		return nil, api.GetError(err)
	}
	as := new(models.AccessSessionDBList)

	err = json.Unmarshal(accessByte, &as)
	log.Print(string(accessByte))
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, api.GetError(err)
	}
	user.Access = as.List()
	p_add := make(map[string]interface{})
	p_add["profile"] = *profile
	user.Additional = p_add
	return user, api.GetError(err)
}

/*
 *	Signin
 */

func SignIn(s_in *models.SignIn) (user *api.UserSession, ap_err *api.ApiError) {
	query := "SELECT u.uuid, u.email, u.confirmed, u.privacy_policy, u.updated, u.created, " +
		"p.uuid, u.uuid, p.first_name, p.last_name, CONCAT(p.first_name, ' ', p.last_name), " +
		"ifnull(p.display_name, 'none'), ifnull(p.gender, 'none'),p.updated, p.created, " +
		"ifnull(av.type, ''), ifnull(av.url,''), ifnull(av.updated, 0), ifnull(av.created, 0), " +
		"JSON_OBJECT(" +
		"a.service, JSON_OBJECT(" +
		"ifnull(m.uuid, 'default'), GROUP_CONCAT(CONCAT(ar.right)))), " +
		"pi.password " +
		"FROM drops_user AS u " +
		"LEFT JOIN password_info AS pi ON pi.drops_user_id = u.id " +
		"LEFT JOIN profile AS p ON p.drops_user_id = u.id " +
		"LEFT JOIN avatar AS av ON av.profile_id = av.id " +
		"LEFT JOIN access AS a ON u.id = a.drops_user_id " +
		"LEFT JOIN access_right AS ar ON ar.access_id = a.id " +
		"LEFT JOIN model AS m ON a.model_id = m.id " +
		"WHERE u.email = ? " +
		"GROUP BY u.id, p.id, av.profile_id, pi.password, a.id, m.id"
	var accessByte []byte
	var password []byte
	user = new(api.UserSession)
	profile := new(models.Profile)
	rows, err := utils.DB.Query(query, s_in.Email)
	if err != nil {
		return nil, api.GetError(err)
	}
	for rows.Next() {
		err = rows.Scan(
			&user.Uuid,
			&user.Email,
			&user.Confirmed,
			&user.PrivacyPolicy,
			&user.Updated,
			&user.Created,
			&profile.Uuid,
			&profile.UserUuid,
			&profile.FirstName,
			&profile.LastName,
			&profile.FullName,
			&profile.DisplayName,
			&profile.Gender,
			&profile.Created,
			&profile.Updated,
			&profile.Avatar.Type,
			&profile.Avatar.Url,
			&profile.Avatar.Updated,
			&profile.Avatar.Created,
			&accessByte,
			&password)
		if err != nil {
			return nil, api.GetError(err)
		}
		//password check
		err = bcrypt.CompareHashAndPassword(password, []byte(s_in.Password))
		if err != nil {
			return nil, api.GetError(api.ErrorPassword)
		}
		as := new(models.AccessSessionDBList)
		err = json.Unmarshal(accessByte, &as)
		log.Print(string(accessByte))
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, api.GetError(err)
		}
		user.Access = as.List()
		p_add := make(map[string]interface{})
		p_add["profile"] = *profile
		user.Additional = p_add
	}
	return user, api.GetError(err)
}
