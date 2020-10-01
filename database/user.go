package database

import (
	"log"
	"strings"
	"time"

	//	"strconv"

	"drops-backend/models"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/api"
	crmm "github.com/Viva-con-Agua/echo-pool/crm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

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
		err = api.ErrorNotFound
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

func PasswordResetToken(n *models.NewToken) (crm_email *crmm.CrmEmailBody, api_err *api.ApiError) {
	ce := n.CrmEmailBody("PASSWORD_RESET")

	//Select drops_user_id
	query := "SELECT du.id, du.uuid, du.country FROM drops_user as du " +
		"JOIN password_info as p ON du.id = p.drops_user_id " +
		"WHERE du.email = ? && du.confirmed = 1 " +
		"LIMIT 1"
	var u_id int
	err := utils.DB.QueryRow(query, n.Email).Scan(&u_id, &ce.CrmData.DropsId, &ce.CrmData.Country)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, api.GetError(utils.ErrorUserNotFound)
		}
		return nil, api.GetError(err)
	}
	//Begin Database Query
	tx, err := utils.DB.Begin()
	if err != nil {
		return nil, api.GetError(err)
	}
	//select password Token id
	query = "DELETE FROM access_token WHERE access_token.drops_user_id = ? && access_token.t_case = 'password'"
	_, err = tx.Exec(query, u_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	// insert new password Token
	token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	_, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, redirect_url, expired, created, drops_user_id) VALUES(?, ?, ?, ?, ?, ?)",
		token,
		"password",
		"none",
		time.Now().Add(time.Hour*24).Unix(),
		time.Now().Unix(),
		u_id,
	)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	ce.Mail.Link = token
	ce.CrmData.Created = time.Now().Unix()
	return ce, api.GetError(tx.Commit())

}

func PasswordReset(pr *models.PasswordReset) (u_uuid *string, api_err *api.ApiError) {
	now := time.Now().Unix()
	var u_id, t_id int64
	query := "SELECT u.id, u.uuid, t.id FROM drops_user AS u " +
		"JOIN access_token AS t ON t.drops_user_id = u.id " +
		"WHERE t.token = ? && t.expired > ?"
	err := utils.DB.QueryRow(query, pr.Token, now).Scan(&u_id, &u_uuid, &t_id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, api.GetError(utils.ErrorUserNotFound)
		}
		return nil, api.GetError(err)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(pr.Password), 10)
	if err != nil {
		return nil, api.GetError(err)
	}
	query = "UPDATE password_info AS pi " +
		"SET pi.password = ?, updated = ?  " +
		"WHERE pi.drops_user_id = ?"
	tx, err := utils.DB.Begin()
	if err != nil {
		return nil, api.GetError(err)
	}
	//update user user
	res, err := tx.Exec(query, password, now, u_id)
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
		err = api.ErrorNotFound
		return nil, api.GetError(err)
	}
	_, err = tx.Exec("DELETE FROM access_token WHERE access_token.id = ?", t_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	return u_uuid, api.GetError(tx.Commit())

}
