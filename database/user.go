package database

import (
	"log"
	"time"

	//	"strconv"

	"drops-backend/models"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/api"
	_ "github.com/go-sql-driver/mysql"
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

func PasswordResetToken(n *models.NewToken) (*string, error) {
	//Begin Database Query
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return nil, err
	}
	//Select drops_user_id
	query := "SELECT du.id FROM drops_user as du " +
		"JOIN password_info as p ON du.id = p.drops_user_id " +
		"WHERE du.email = ? && du.confirmed = 1 &&  " +
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
		return nil, api.ErrorNotFound
	}

	//select password Token id
	query = "SELECT id FROM access_token WHERE access_token.id = ? && access_token.t_case = 'password'"
	rows, err = utils.DB.Query(query, id)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	var access_id int
	// Delete Token if there is one
	for rows.Next() {
		err = rows.Scan(&access_id)
		if err != nil {
			query = "DELETE FROM access_token WHERE id=?"
			_, err := tx.Exec(query, id)
			if err != nil {
				log.Print("Database Error", err)
				return nil, err
			}

		}
	}
	// insert new password Token
	access_token, err := utils.RandomBase64(32)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}
	_, err = tx.Exec(
		"INSERT INTO access_token (token, t_case, expired, created, drops_user_id) VALUES(?, ?, ?, ?, ?)",
		access_token,
		"password",
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
