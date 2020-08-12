package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/**
 * select a list of models.Roles
 */
func GetRolesDefaultList(page *models.Page, sort string, filter *models.FilterRole) (Roles []models.CrewRoleExtended, err error) {

	rolesQuery := "SELECT cr.uuid, cr.name, cr.description, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', p.uuid, 'name', p.name, 'description', p.description)) " +
		", ']') " +
		"FROM crew_role AS cr " +
		"LEFT JOIN crew_role_has_permission AS crp ON cr.id = crp.crew_role_id " +
		"LEFT JOIN permission AS p ON crp.permission_id = p.id " +
		"WHERE cr.name LIKE ? " +
		"GROUP BY cr.id " +
		sort + " " +
		"LIMIT ?, ?"

	// Execute the Query
	rows, err := utils.DB.Query(rolesQuery, filter.Name, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	var permissionByte []byte

	// convert each row
	for rows.Next() {

		//create models
		role := new(models.CrewRoleExtended)
		permission := new(models.PermissionList)

		// scan row and fill Profile
		err = rows.Scan(&role.Uuid, &role.Name, &role.Description, &permissionByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add addresses
		err = json.Unmarshal(permissionByte, &permission)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// add addresses to user
		if (*permission)[0].Uuid != "" {
			role.Permissions = *permission.Distinct()
		} else {
			role.Permissions = nil
		}

		Roles = append(Roles, *role)
	}

	return Roles, err
}

/**
 * Create Role
 */
func CreateRole(Role *models.CrewRoleCreate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing profile
	rows, err := tx.Query("SELECT id FROM crew_role WHERE name = ?", Role.Name)
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

	// Insert Profile
	uuid := uuid.New()
	_, err = tx.Exec("INSERT INTO crew_role (uuid, name, description, updated, created) VALUES "+
		"(?, ?, ?, ?, ?)",
		uuid, Role.Name, Role.Description, time.Now().Unix(), time.Now().Unix())
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
func GetRole(search string) (Roles []models.CrewRoleExtended, err error) {

	rolesQuery := "SELECT cr.uuid, cr.name, cr.description, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', p.uuid, 'name', p.name, 'description', p.description)) " +
		", ']') " +
		"FROM crew_role AS cr " +
		"LEFT JOIN crew_role_has_permission AS crp ON cr.id = crp.crew_role_id " +
		"LEFT JOIN permission AS p ON crp.permission_id = p.id " +
		"WHERE cr.uuid LIKE ? " +
		"GROUP BY cr.id"

	log.Print(rolesQuery)

	// Execute the Query
	rows, err := utils.DB.Query(rolesQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	var permissionByte []byte

	// convert each row
	for rows.Next() {

		//create models
		role := new(models.CrewRoleExtended)
		permission := new(models.PermissionList)

		// scan row and fill Profile
		err = rows.Scan(&role.Uuid, &role.Name, &role.Description, &permissionByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add addresses
		err = json.Unmarshal(permissionByte, &permission)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// add addresses to user
		if (*permission)[0].Uuid != "" {
			role.Permissions = *permission.Distinct()
		} else {
			role.Permissions = nil
		}

		Roles = append(Roles, *role)
	}

	return Roles, err
}

/**
 * update Profile
 */
func UpdateRole(Role *models.CrewRoleUpdate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM crew_role WHERE uuid = ?", Role.Uuid)

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

	//update role
	_, err = tx.Exec("UPDATE crew_role SET name = ?, description = ?, updated = ? WHERE id = ?", Role.Name, Role.Description, time.Now().Unix(), id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

func DeleteRole(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM crew_role WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		log.Print("Database Error: ", err)
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

	// Delete assigned roles
	_, err = tx.Exec("DELETE FROM profile_has_crew_has_crew_role WHERE crew_role_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete role permissions
	_, err = tx.Exec("DELETE FROM crew_role_has_permission WHERE crew_role_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete role itself
	_, err = tx.Exec("DELETE FROM crew_role WHERE id = ?", id)
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
func AssignRole(Assignment *models.AssignCrewRole) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing profile
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", Assignment.ProfileId)
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

	// Check for existing relation
	rows, err = tx.Query("SELECT id FROM profile_has_crew WHERE profile_id = ? AND crew_id = ?", profileId, crewId)
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

	// Check for existing role
	rows, err = tx.Query("SELECT id FROM crew_role WHERE uuid = ?", Assignment.RoleId)
	if err != nil {
		log.Print("Database Error", err)
		return err
	}

	var roleId int
	for rows.Next() {
		err = rows.Scan(&roleId)
		if err != nil {
			log.Print("Database Error: ", err)
			return err
		}
	}
	//if roleId == 0 return NotFound
	if roleId == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Insert Profile
	uuid := uuid.New()
	_, err = tx.Exec("INSERT INTO profile_has_crew_has_crew_role (uuid, profile_has_crew_id, crew_role_id, updated, created) VALUES "+
		"(?, ?, ?, ?, ?)",
		uuid, id, roleId, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

func RemoveRole(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT count(*) FROM profile_has_crew_has_crew_role WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		log.Print("Database Error: ", err)
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
	//if id == 0 return NotFound
	if cnt == 0 {
		err = utils.ErrorNotFound
		return err
	}

	// Delete assigned roles
	_, err = tx.Exec("DELETE FROM profile_has_crew_has_crew_role WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}
