package database

import (
	"log"

	"../models"
	"../utils"
	"github.com/google/uuid"
)

/**
 * insert role into database
 */
func PostRole(r *models.RoleCreate) (err error) {
	// Create uuid
	Uuid, err := uuid.NewRandom()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	// begin database query and handle error
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	// insert role
	_, err = tx.Exec("INSERT INTO Role (uuid, name, pillar) VALUES(?, ?, ?)", Uuid.String(), r.Name, r.Pillar)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}

/**
 * select a list of models.Roles
 */
func GetRolesList() (roles []models.Role, err error) {
	// Execute the Query
	rows, err := utils.DB.Query("SELECT * FROM Role")
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	// define variable for each column
	var id int
	var uuid, name, pillar string
	// convert each row
	for rows.Next() {
		//scan row
		err = rows.Scan(&id, &uuid, &name, &pillar)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// fill models.Role
		role := new(models.Role)
		role.Uuid = uuid
		role.Name = name
		role.Pillar = pillar
		roles = append(roles, *role)
	}
	return roles, err
}

/**
 * select roles for an given Supporter_id
 */
func GetRolesByProfileId(Supporter_id int) (roles []models.Role, err error) {
	// Execute the Query
	rows, err := utils.DB.Query("SELECT * FROM Supporter_has_Role LEFT JOIN Role ON Supporter_has_Role.Role_Id = Role.id WHERE Supporter_has_Role.Supporter_id = ?", Supporter_id)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	// define variable for each column
	var id int
	var uuid, name, pillar string
	// convert each row
	for rows.Next() {
		//scan row
		err = rows.Scan(&id, &id, &id, &uuid, &name, &pillar)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// fill models.Role
		role := new(models.Role)
		role.Uuid = uuid
		role.Name = name
		role.Pillar = pillar
		roles = append(roles, *role)
	}
	return roles, err
}
