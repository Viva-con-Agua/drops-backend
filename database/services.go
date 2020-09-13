package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"strings"

	"github.com/google/uuid"
)

/*
 *
 */
func ServiceCreate(service_create *models.ServiceCreate) (service *models.Service, err error) {
	//initial database
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Error: database.ServiceInsert Step_1 ### ", err)
		return nil, err
	}
	query := "INSERT INTO service (uuid, name, description, created, updated) " +
		"VALUES(?, ?, ?, ?, ?)"
	service = service_create.Service()
	res, err := tx.Exec(
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
			return nil, utils.ErrorConflict
		}
		log.Print("Error: database.ServiceInsert Step_2 ### ", err)
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Error: database.ServiceInsert Step_3 ### ", err)
		return nil, err
	}
	uuid_model := uuid.New().String()
	query = "INSERT INTO model (uuid, name, type, created, service_id) " +
		"VALUES(?, ?, ?, ?, ?)"
	_, err = tx.Exec(
		query,
		uuid_model,
		"default",
		"control",
		id,
	)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, utils.ErrorConflict
		}
		log.Print("Error: database.ServiceInsert Step_4 ### ", err)
		return nil, err
	}
	return service, tx.Commit()
}

func ServiceList() (service_list []models.Service, err error) {
	query := "SELECT s.id, s.uuid, s.name, s.description, s.created, s.updated " +
		"FROM service AS s " +
		"GROUP BY s.id"

	rows, err := utils.DB.Query(query)
	if err != nil {
		log.Print("Error: database.ServiceList Step_1 ### ", err)
		return nil, err
	}
	var id int64
	service := new(models.Service)
	for rows.Next() {
		err = rows.Scan(
			&id,
			&service.Uuid,
			&service.Name,
			&service.Description,
			&service.Created,
			&service.Updated,
		)
		if err != nil {
			log.Print("Error: database.ServiceList Step_2 ### ", err)
			return nil, err
		}
		service_list = append(service_list, *service)
	}
	return service_list, err
}
