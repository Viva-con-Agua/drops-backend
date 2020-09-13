package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"strings"
)

func AddEssentialModels() (err error) {
	//initial database
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Error: database.ServiceInsert Step_1 ### ", err)
		return err
	}
	service_create := models.ServiceCreate{Name: "drops-backend", Description: "service for handling all vca_user"}
	query := "INSERT INTO service (uuid, name, description, created, updated) " +
		"VALUES(?, ?, ?, ?, ?)"
	service := service_create.Service()
	_, err = tx.Exec(
		query,
		service.Uuid,
		service.Name,
		service.Description,
		service.Created,
		service.Updated,
	)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrorConflict
		}
		log.Print("Error: database.ServiceInsert Step_2 ### ", err)
		return err
	}
	/*
		id, err := res.LastInsertId()
		if err != nil {
			log.Print("Error: database.ServiceInsert Step_3 ### ", err)
			return err
		}
		uuid_model := uuid.New().String()
		query = "INSERT INTO model (uuid, name, type, created, service_id) " +
			"VALUES(?, ?, ?, ?, ?)"
		_, err = tx.Exec(
			query,
			uuid_model,
			"default",
			"control",
			time.Now().Unix(),
			id,
		)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "Duplicate entry") {
				return utils.ErrorConflict
			}
			log.Print("Error: database.ServiceInsert Step_4 ### ", err)
			return err
		}*/

	return tx.Commit()

}
