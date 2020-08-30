package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func AdminCreate(c *models.SignUpData) (*string, error) {
	// Create uuid
	Uuid := uuid.New()

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
		c.SignUpUser.Email,
		1,
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
	password, err := bcrypt.GenerateFromPassword([]byte(c.SignUpUser.Password), 10)
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
		c.SignUpUser.Email,
		c.SignUpUser.LastName,
		time.Now().Unix(),
		time.Now().Unix(),
		id,
	)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return nil, err
	}

	result := Uuid.String()
	// insert profile
	return &result, tx.Commit()
}
