package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"

	"github.com/google/uuid"
)

/**
 * join user and role entry via User_has_Role table
 */
func AccessInsert(assign *models.AccessCreate) (err error) {
	// select user_id from database
	rows, err := utils.DB.Query("SELECT id FROM vca_user WHERE uuid = ?", assign.Assign)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}
	// select user_id from rows
	var userId int
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	// select role_id from database
	rows2, err := utils.DB.Query("SELECT id FROM role WHERE uuid = ?", assign.RoleUuid)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}
	//select user_id from rows
	var roleId int
	for rows2.Next() {
		err = rows2.Scan(&roleId)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	// begin database query and handle error
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	// Create uuid
	Uuid, err := uuid.NewRandom()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	res, err := tx.Exec("INSERT INTO access_user (uuid, pool_user_id, role_Id) VALUES(?, ?, ?)", Uuid, userId, roleId)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	if assign.ModelUuid != "" && assign.ModelName != "" {
		// get user id via LastInsertId
		id, err := res.LastInsertId()
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
		_, err = tx.Exec("INSERT INTO model (uuid, name, access_user_id) VALUES(?, ?, ?)", assign.ModelUuid, assign.ModelName, id)
		if err != nil {
			tx.Rollback()
			log.Print("Database Error: ", err)
			return err
		}
	}
	return tx.Commit()
}

/**
 * DELETE /users/access
 */
func AccessDelete(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	rows, err := tx.Query("SELECT id FROM access_user WHERE uuid = ?", deleteBody.Uuid)
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
	_, err = tx.Exec("DELETE FROM access_user WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
