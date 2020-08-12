package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

/**
 * Create Profile
 */
func CreateProfile(Profile *models.ProfileCreate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// Check for existing profile
	rows, err := tx.Query("SELECT id FROM profile WHERE email = ?", Profile.Email)
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
	uuid := uuid.New()
	res, err := tx.Exec("INSERT INTO profile (uuid, firstname, lastname, email, mobile, birthdate, gender, updated, created) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		uuid, Profile.FirstName, Profile.LastName, Profile.Email, Profile.Mobile, Profile.Birthdate, Profile.Gender, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Create default newsletter settings
	profileId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Insert Newsletter setting
	_, err = tx.Exec("INSERT INTO newsletter (profile_id, setting) VALUES "+
		"(?, ?)",
		profileId, "none")
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

/**
 * select Profile by uuid
 */
func GetProfile(search string) (Profiles []models.ProfileExtended, err error) {

	// execute the query
	profileQuery := "SELECT profile.uuid, profile.email, profile.firstname, profile.lastname,  CONCAT(profile.firstname, ' ', profile.lastname) AS fullname, profile.mobile, profile.birthdate, profile.gender, profile.updated, profile.created, " +
		"IFNULL(avatar.url, ''), IFNULL(avatar.type, ''), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', address.uuid, 'primary', profile_has_address.primary, 'street', address.street, 'additional', address.additional, 'zip', address.zip, 'city', address.city, 'country', address.country, 'google_id', address.google_id, 'updated', address.updated, 'created', address.created)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', crew.uuid, 'primary', nvm_state.primary_crew, 'name', crew.name, 'email', crew.email, 'abbreviation', crew.abbreviation)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('crew_uuid', crew.uuid, 'crew_name', crew.name, 'role', crew_role.name)) " +
		", ']') " +
		"FROM profile " +
		"LEFT JOIN avatar ON avatar.profile_id = profile.id " +
		"LEFT JOIN profile_has_crew ON profile_has_crew.profile_id = profile.id " +
		"LEFT JOIN profile_has_crew_has_crew_role ON profile_has_crew_has_crew_role.profile_has_crew_id = profile_has_crew.id " +
		"LEFT JOIN crew_role ON crew_role.id = profile_has_crew_has_crew_role.crew_role_id " +
		"LEFT JOIN nvm_state ON nvm_state.profile_has_crew_id = profile_has_crew.id " +
		"LEFT JOIN crew ON profile_has_crew.crew_id = crew.id  " +
		"LEFT JOIN profile_has_address ON profile.id = profile_has_address.profile_id " +
		"LEFT JOIN address ON address.id = profile_has_address.address_id " +
		"WHERE profile.uuid = ? " +
		"GROUP BY profile.id "

	rows, err := utils.DB.Query(profileQuery, search)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	var addressByte []byte
	var crewsByte []byte
	var rolesByte []byte

	// convert each row
	for rows.Next() {

		//create Profile and corresponding models
		profile := new(models.ProfileExtended)
		avatar := new(models.Avatar)
		crews := new(models.CrewList)
		address := new(models.AddressList)
		roles := new(models.ProfileRoleList)

		//scan row and fill Profile
		err = rows.Scan(&profile.Uuid, &profile.Email, &profile.FirstName, &profile.LastName, &profile.FullName, &profile.Mobile, &profile.Birthdate, &profile.Gender, &profile.Updated, &profile.Created, &avatar.Url, &avatar.Type, &addressByte, &crewsByte, &rolesByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		profile.Avatar = *avatar

		// Add addresses
		err = json.Unmarshal(addressByte, &address)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add roles to user
		if (*address)[0].Uuid != "" {
			profile.Addresses = *address.Distinct()
		} else {
			profile.Addresses = nil
		}

		// Add crews
		err = json.Unmarshal(crewsByte, &crews)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// add crews to user
		if (*crews)[0].Uuid != "" {
			profile.Crews = *crews.Distinct()
		} else {
			profile.Crews = nil
		}

		// Add roles
		err = json.Unmarshal(rolesByte, &roles)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}
		// add roles to user
		if (*roles)[0].Uuid != "" {
			profile.Roles = *roles.Distinct().NotEmpty()
		} else {
			profile.Roles = nil
		}

		Profiles = append(Profiles, *profile)
	}
	return Profiles, err
}

/**
 * update Profile
 */
func UpdateProfile(Profile *models.ProfileUpdate) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", Profile.Uuid)

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
	_, err = tx.Exec("UPDATE profile SET firstname = ?, lastname = ?, mobile = ?, birthdate = ?, gender = ?, updated = ? WHERE id = ?", Profile.FirstName, Profile.LastName, Profile.Mobile, Profile.Birthdate, Profile.Gender, time.Now().Unix(), id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

func DeleteProfile(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}

	// select id
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", deleteBody.Uuid)
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

	// Delete addresses
	_, err = tx.Exec("DELETE FROM newsletter WHERE profile_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete avatar
	_, err = tx.Exec("DELETE FROM avatar WHERE profile_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Get address ids

	// Delete address relation
	_, err = tx.Exec("DELETE FROM profile_has_address WHERE profile_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned crew_roles of profile
	_, err = tx.Exec("DELETE FROM profile_has_crew_has_crew_role WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE profile_id = ?)", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned crew_roles of profiles
	_, err = tx.Exec("DELETE FROM nvm_state WHERE profile_has_crew_id IN (SELECT id FROM profile_has_crew WHERE profile_id = ?)", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete assigned crews profile
	_, err = tx.Exec("DELETE FROM profile_has_crew WHERE profile_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	// Delete profile itself
	_, err = tx.Exec("DELETE FROM profile WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * update Profile
 */
func SetNewsletter(Setting *models.ProfileNewsletter) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", Setting.Uuid)

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
	_, err = tx.Exec("UPDATE newsletter SET setting = ? WHERE profile_id = ?", Setting.Setting, id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()

}

/**
 * update Profile
 */
func UpdateAvatar(Avatar *models.AvatarUpdate, File multipart.File) (err error) {
	// sgl begin
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	//slect id
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", Avatar.Uuid)

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
	_, err = tx.Exec("UPDATE avatar SET type = ?, url = ?, updated= ? WHERE id = ?", Avatar.Type, Avatar.Url, time.Now().Unix(), id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()

}

func DeleteAvatar(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	// select id
	rows, err := tx.Query("SELECT id FROM profile WHERE uuid = ?", deleteBody.Uuid)
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

	// Delete avatar
	_, err = tx.Exec("DELETE FROM avatar WHERE profile_id = ?", id)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}

	return tx.Commit()
}
