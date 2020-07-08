package database

import (
	"../models"
	"../utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"time"
)

/**
 * select Crew
 */
func GetCrew(search string) (Crews []models.CrewExtended, err error) {
	// execute the query
	CrewQuery := "SELECT Crew.id, Crew.uuid, Profile.email, Profile.first_name, Profile.last_name, Crew.updated, Crew.created " +
		"FROM Crew LEFT JOIN Profile ON Crew.id = Profile.Crew_id " +
		"LEFT JOIN Crew_has_Role ON Crew.id = Crew_has_Role.Crew_id " +
		"LEFT JOIN Role ON Crew_has_Role.Role_Id = Role.id " +
		"WHERE Crew.uuid = ? " +
		"GROUP BY Crew.id "
	rows, err := utils.DB.Query(CrewQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	//define variable for Crew database id
	var id int
	Crew := new(models.CrewExtended)
	for rows.Next() {
		//scan row
		err = rows.Scan(&id, &Crew.Uuid, &Crew.Name, &Crew.Updated, &Crew.Created)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		log.Print(id)
		// get roles by id
		var cities []models.City
		cities, err = GetCitiesByCrewId(id)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		//join roles to Crew
		Crew.Cities = cities
	}
	if id == 0 {
		return nil, err
	}
	Crews = append(Crews, *Crew)
	return Crews, err
}

/**
 * select list of Crew
 */
func GetCrewList(page *models.Page, sort string, filter *models.FilterCrew) (Crews []models.CrewExtended, err error) {
	// execute the query
	CrewQuery := "SELECT u.id, u.uuid, p.email, p.first_name, p.last_name, u.updated, u.created " +
		"FROM Crew AS u LEFT JOIN Profile AS p ON u.id = p.Crew_id " +
		"LEFT JOIN Crew_has_Role ON u.id = Crew_has_Role.Crew_id " +
		"LEFT JOIN Role ON Crew_has_Role.Role_Id = Role.id " +
		"WHERE p.email LIKE ? " +
		"GROUP BY u.id " +
		sort + " " +
		"LIMIT ?, ?"
	rows, err := utils.DB.Query(CrewQuery, filter.Name, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	// variable for Crew database id
	var id int
	// convert each row
	for rows.Next() {
		//create Crew
		Crew := new(models.CrewExtended)
		//scan row and fill Crew
		err = rows.Scan(&id, &Crew.Uuid, &Crew.Name, &Crew.Updated, &Crew.Created)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// get cities by id
		var cities []models.City
		cities, err = GetCitiesByCrewId(id)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// join cities to Crew
		Crew.Cities = cities
		// append to list of Crew
		Crews = append(Crews, *Crew)
	}
	return Crews, err
}

/**
 * update Crew
 */
func UpdateCrew(Crew *models.CrewUpdate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM Crew WHERE uuid = ?", Crew.Uuid)
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

	//update Crew Crew
	_, err = tx.Exec("UPDATE Crew SET updated = ? WHERE id = ?", time.Now().Unix(), id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	//update Cities
	_, err = tx.Exec("UPDATE Profile SET first_name = ?, last_name = ? WHERE Crew_id = ?", Crew.Name, Crew.Name, id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * Create Crew
 */
func CreateCrew(Crew *models.CrewCreate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	//Insert Crew Crew
	id := uuid.New()
	_, err = tx.Exec("INSERT INTO Crew SET id = ?, name = ?, created = ?, updated = ?", id, Crew.Name, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	//TODO Inser Cities
	_, err = tx.Exec("UPDATE Profile SET first_name = ?, last_name = ? WHERE Crew_id = ?", Crew.Name, Crew.Name, id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * Delete Crew
 */
func DeleteCrew(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM Crew WHERE uuid = ?", deleteBody.Uuid)
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

	//update Crew Crew
	_, err = tx.Exec("DELETE FROM Crew WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}
