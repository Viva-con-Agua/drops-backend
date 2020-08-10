package database

import (
	"../models"
	"../utils"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"time"
)

/**
 * select list of Crew
 */
func GetCrewList(page *models.Page, sort string, filter *models.FilterCrew) (Crews []models.CrewExtended, err error) {
	// execute the query
	CrewQuery := "SELECT c.uuid, c.name, c.email, c.abbreviation, c.updated, c.created, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('key', cm.meta_key, 'value', cm.meta_value)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', ci.uuid, 'city', ci.city, 'country', ci.country, 'google_id', ci.google_id)) " +
		", ']') " +
		"FROM crew AS c LEFT JOIN crew_meta AS cm ON c.id = cm.crew_id " +
		"LEFT JOIN city AS ci ON c.id = ci.crew_id " +
		"WHERE c.name LIKE ? " +
		"GROUP BY c.id " +
		sort + " " +
		"LIMIT ?, ?"
	rows, err := utils.DB.Query(CrewQuery, filter.Name, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	// variable for Crew database id
	var cityByte []byte
	var metaByte []byte

	// convert each row
	for rows.Next() {
		//create Crew
		crew := new(models.CrewExtended)
		meta := new(models.DictList)
		cities := new(models.CityList)

		//scan row and fill Crew
		err = rows.Scan(&crew.Uuid, &crew.Name, &crew.Email, &crew.Abbreviation, &crew.Updated, &crew.Created, &metaByte, &cityByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add meta
		err = json.Unmarshal(metaByte, &meta)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		log.Print((*meta)[0].Key)

		// add addresses to user
		if (*meta)[0].Key != "" {
			crew.CrewMeta = *meta
		} else {
			crew.CrewMeta = nil
		}

		// Add cities
		err = json.Unmarshal(cityByte, &cities)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add addresses to user
		if (*cities)[0].Uuid != "" {
			crew.Cities = *cities
		} else {
			crew.Cities = nil
		}

		// append to list of Crew
		Crews = append(Crews, *crew)
	}
	return Crews, err
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

	// Check for existing profile
	rows, err := tx.Query("SELECT id FROM crew WHERE name = ?", Crew.Name)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFound
	if id != 0 {
		err = utils.ErrorConflict
		return err
	}

	// Insert Profile
	crewUuid := uuid.New()
	res, err := tx.Exec("INSERT INTO crew (uuid, name, email, abbreviation, updated, created) VALUES "+
		"(?, ?, ?, ?, ?, ?)",
		crewUuid, Crew.Name, Crew.Email, Crew.Abbreviation, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	crewId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Insert Cities
	for _, c := range Crew.Cities {

		cityUuid := uuid.New()
		res, err = tx.Exec("INSERT INTO city (uuid, city, country, google_id, crew_id, updated, created) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?)",
			cityUuid, c.City, c.Country, c.GoogleId, crewId, time.Now().Unix(), time.Now().Unix())
		if err != nil {
			tx.Rollback()
			log.Print("Database Error: ", err)
			return err
		}
	}

	return tx.Commit()

}

/**
 * select Crew
 */
func GetCrew(search string) (Crews []models.CrewFull, err error) {
	// execute the query
	CrewQuery := "SELECT c.uuid, c.name, c.email, c.abbreviation, c.updated, c.created, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('key', cm.meta_key, 'value', cm.meta_value)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', ci.uuid, 'city', ci.city, 'country', ci.country, 'google_id', ci.google_id)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('id', cr.id, 'name', cr.name, 'description', cr.description, 'first_name', p.firstname, 'last_name', p.lastname)) " +
		", ']') " +
		"FROM crew AS c LEFT JOIN crew_meta AS cm ON c.id = cm.crew_id " +
		"LEFT JOIN city AS ci ON c.id = ci.crew_id " +
		"LEFT JOIN profile_has_crew AS pc ON pc.crew_id = c.id " +
		"LEFT JOIN profile AS p ON p.id = pc.profile_id " +
		"LEFT JOIN profile_has_crew_has_crew_role AS pcr ON pcr.profile_has_crew_id = pc.id " +
		"LEFT JOIN crew_role AS cr ON cr.id = pcr.crew_role_id " +
		"WHERE c.uuid = ? " +
		"GROUP BY c.id "
	rows, err := utils.DB.Query(CrewQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}
	// variable for Crew database id
	var cityByte []byte
	var metaByte []byte
	var roleByte []byte

	// convert each row
	for rows.Next() {
		//create Crew
		crew := new(models.CrewFull)
		meta := new(models.DictList)
		cities := new(models.CityList)
		roles := new(models.CrewRoleProfileList)

		//scan row and fill Crew
		err = rows.Scan(&crew.Uuid, &crew.Name, &crew.Email, &crew.Abbreviation, &crew.Updated, &crew.Created, &metaByte, &cityByte, &roleByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add meta
		err = json.Unmarshal(metaByte, &meta)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add meta to crew
		if (*meta)[0].Key != "" {
			crew.CrewMeta = *meta.Distinct()
		} else {
			crew.CrewMeta = nil
		}

		// Add cities
		err = json.Unmarshal(cityByte, &cities)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add city to crew
		if (*cities)[0].Uuid != "" {
			crew.Cities = *cities.Distinct()
		} else {
			crew.Cities = nil
		}

		// Add roles
		err = json.Unmarshal(roleByte, &roles)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add roles to crew
		if (*roles)[0].Name != "" {
			crew.Roles = *roles
		} else {
			crew.Roles = nil
		}

		// append to list of Crew
		Crews = append(Crews, *crew)
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
	rows, err := tx.Query("SELECT id FROM crew WHERE uuid = ?", Crew.Uuid)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

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

	//update profile
	_, err = tx.Exec("UPDATE crew SET name = ?, email = ?, abbreviation = ?, updated = ? WHERE id = ?", Crew.Name, Crew.Email, Crew.Abbreviation, time.Now().Unix(), id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	//update Cities

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
	// select ids
	rows, err := tx.Query("SELECT id FROM crew WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

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

	// Delete crew_meta
	_, err = tx.Exec("DELETE FROM crew_meta WHERE crew_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete cities
	_, err = tx.Exec("DELETE FROM city WHERE crew_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned crew_roles of profiles
	_, err = tx.Exec("DELETE FROM profile_has_crew_has_crew_role WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE crew_id = ?)", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned crew_roles of profiles
	_, err = tx.Exec("DELETE FROM nvm_state WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE crew_id = ?)", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned profiles
	_, err = tx.Exec("DELETE FROM profile_has_crew WHERE crew_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	//update Crew Crew
	_, err = tx.Exec("DELETE FROM crew WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * Assign role to profile crew relation
 */
func AssignCrew(Assignment *models.AssignCrew) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing relation
	rows, err := tx.Query("SELECT pc.id FROM profile_has_crew AS pc "+
		"INNER JOIN profile AS p ON p.id = pc.profile_id AND p.uuid = ? "+
		"INNER JOIN crew AS c ON c.id = pc.crew_id AND c.uuid = ? ",
		Assignment.ProfileId, Assignment.CrewId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id != 0 return Conflict
	if id != 0 {
		err = utils.ErrorConflict
		return err
	}

	// Check for existing profile
	rows, err = tx.Query("SELECT id FROM profile WHERE uuid = ?", Assignment.ProfileId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var profileId int
	for rows.Next() {
		err = rows.Scan(&profileId)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFound
	if profileId == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Check for existing crew
	rows, err = tx.Query("SELECT id FROM crew WHERE uuid = ?", Assignment.CrewId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var crewId int
	for rows.Next() {
		err = rows.Scan(&crewId)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFound
	if crewId == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Insert Relation
	res, err := tx.Exec("INSERT INTO profile_has_crew (profile_id, crew_id, updated, created) VALUES "+
		"(?, ?, ?, ?)",
		profileId, crewId, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	var relId int64
	relId, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	if Assignment.Primary {

		//update nvm_states
		_, err = tx.Exec("UPDATE nvm_state SET primary_crew = false WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE profile_id = ?)", profileId)
		if err != nil {
			tx.Rollback()
			log.Print("Database Error: ", err)
			return err
		}

	}

	active_state := "inactive"
	if Assignment.Active {
		active_state = "requested"
	}

	// Insert Relation
	_, err = tx.Exec("INSERT INTO nvm_state (profile_has_crew_id, active_state, primary_crew) VALUES "+
		"(?, ?, ?)",
		relId, active_state, Assignment.Primary)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

/**
 * Remove role from profile crew relation
 */
func RemoveCrew(DeleteBody *models.RemoveCrew) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing relation
	rows, err := tx.Query("SELECT pc.id FROM profile_has_crew AS pc "+
		"INNER JOIN profile AS p ON p.id = pc.profile_id AND p.uuid = ? "+
		"INNER JOIN crew AS c ON c.id = pc.crew_id AND c.uuid = ? ",
		DeleteBody.ProfileId, DeleteBody.CrewId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFount
	if id == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Delete assigned crew_roles of profiles
	_, err = tx.Exec("DELETE FROM profile_has_crew_has_crew_role WHERE profile_has_crew_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned nvm_states of profiles
	_, err = tx.Exec("DELETE FROM nvm_state WHERE profile_has_crew_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned profiles
	_, err = tx.Exec("DELETE FROM profile_has_crew WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

/**
 * update Active State
 */
func ActiveStateChange(ActiveState *models.ActiveState, State string) (err error) {

	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing relation
	rows, err := tx.Query("SELECT pc.id FROM profile_has_crew AS pc "+
		"INNER JOIN profile AS p ON p.id = pc.profile_id AND p.uuid = ? "+
		"INNER JOIN crew AS c ON c.id = pc.crew_id AND c.uuid = ? ",
		ActiveState.ProfileId, ActiveState.CrewId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFount
	if id == 0 {
		err = utils.ErrorNotFound
		return err
	}

	//update profile
	_, err = tx.Exec("UPDATE nvm_state SET active_state = ? WHERE profile_has_crew_id = ?", State, id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	//update Cities

	return tx.Commit()

}

/**
 * update Crew
 */
func NVMStateChange(NVMState *models.NVMState, State string) (err error) {

	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing relation
	rows, err := tx.Query("SELECT pc.id FROM profile_has_crew AS pc "+
		"INNER JOIN profile AS p ON p.id = pc.profile_id AND p.uuid = ? "+
		"INNER JOIN crew AS c ON c.id = pc.crew_id AND c.uuid = ? ",
		NVMState.ProfileId, NVMState.CrewId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFount
	if id == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Check for existing precondition
	rows, err = tx.Query("SELECT count(*) FROM nvm_state "+
		"WHERE profile_has_crew_id = ? AND active_state = 'active' AND primary_crew = 1",
		id)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var cnt int
	for rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if id == 0 return NotFount
	if cnt == 0 {
		err = utils.ErrorNotFound
		return err
	}

	var nvmState int64
	if State == "active" {
		nvmState = time.Now().AddDate(1, 0, 0).Unix()

		//update nvm_states
		_, err = tx.Exec("UPDATE nvm_state SET nvm_expiration_date = NULL WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE profile_id = (SELECT id FROM profile WHERE uuid = ?))", NVMState.ProfileId)
		if err != nil {
			tx.Rollback()
			log.Print("Database Error: ", err)
			return err
		}
	}

	//update profile
	_, err = tx.Exec("UPDATE nvm_state SET nvm_expiration_date = ? WHERE profile_has_crew_id = ?", nvmState, id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	//update Cities

	return tx.Commit()

}
