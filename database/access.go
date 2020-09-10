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

func AccessInsertDefault(access_default *models.AccessDefault) (err error) {
	rows, err := utils.DB.Query("SELECT id FROM vca_user WHERE uuid = ?", access_default.UserUuid)
	if err != nil {
		log.Print("Error: database.AccessInsertDefault Step_1 ### ", err)
		return err
	}
	// select user_id from rows
	var userId int
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Print("Error: database.AccessInsertDefault Step_2 ### ", err)
			return err
		}
	}
	rows, err = utils.DB.Query("SELECT id FROM model As m JOIN service AS s WHERE m.name = ? AND s.name = ?", "default", access_default.ServiceName)
	if err != nil {
		log.Print("Error: database.AccessInsertDefault Step", err)
		return err
	}
	// select user_id from rows
	var modelId int
	for rows.Next() {
		err = rows.Scan(&modelId)
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
	_, err = tx.Exec("INSERT INTO access_user (uuid, name, created, vca_user_id, model_id) VALUES(?, ?, ?, ?, ?)",
		Uuid,
		access_default.AccessType,
		userId,
		modelId,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
func AccessInsert(assign *models.AccessCreate) (err error) {
	// select user_id from database
	rows, err := utils.DB.Query("SELECT id FROM vca_user WHERE uuid = ?", assign.Assign)
	if err != nil {
		log.Print("Error: database.AccessInsert Step_1 ### ", err)
		return err
	}
	// select user_id from rows
	var userId int
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Print("Error: database.AccessInsert Step_2 ### ", err)
			return err
		}
	}
	rows, err = utils.DB.Query("SELECT id FROM model WHERE uuid = ?", assign.ModelUuid)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}
	// select user_id from rows
	var modelId int
	for rows.Next() {
		err = rows.Scan(&modelId)
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
	_, err = tx.Exec("INSERT INTO access_user (uuid, name, created, vca_user_id, model_id) VALUES(?, ?, ?, ?, ?)",
		Uuid,
		assign.Name,
		userId,
		modelId,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
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
