package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

func ModelInsert(m *models.ModelCreate) (err error) {

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
	}
	log.Print(m)
	//select user id
	query := "SELECT vca_user.id FROM vca_user " +
		"WHERE vca_user.uuid = ? " +
		"LIMIT 1"

	rows, err := tx.Query(query, m.Owner)
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	var user_id int
	for rows.Next() {
		err = rows.Scan(&user_id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}

	log.Print(user_id)
	query = "INSERT INTO model (uuid, name, service_name, type, created) " +
		"VALUES(?, ?, ?, ?, ?)"

	res, err := tx.Exec(query, m.Uuid, m.Name, m.ServiceName, m.Type, time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	model_id, err := res.LastInsertId()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	log.Print(model_id)
	Uuid, err := uuid.NewRandom()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	query = "INSERT INTO access_user (uuid, name, created, vca_user_id, model_id) VALUES(?, ?, ?, ?, ?)"
	_, err = tx.Exec(query, Uuid, "created", time.Now().Unix(), user_id, model_id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}

func ModelDelete(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM model WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
