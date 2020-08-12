package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"encoding/json"
	"log"
)

/**
 * select list of Profile ids
 */
func GetProfileIdList(page *models.Page, sort string, filter *models.FilterProfile) (Profiles []models.ProfileId, err error) {

	// execute the query
	profileQuery := "SELECT profile.uuid " +
		"FROM profile " +
		"WHERE profile.email LIKE ? " +
		"GROUP BY profile.id " +
		sort + " " +
		"LIMIT ?, ?"

	// check for database error
	rows, err := utils.DB.Query(profileQuery, filter.Email, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	// convert each row
	for rows.Next() {

		// create ProfileId
		profileId := new(models.ProfileId)

		// scan row and fill ProfileId
		err = rows.Scan(&profileId.Uuid)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		Profiles = append(Profiles, *profileId)
	}
	return Profiles, err
}

/**
 * select list of Profile with addresses
 */
func GetProfileDefaultList(page *models.Page, sort string, filter *models.FilterProfile) (Profiles []models.ProfileDefault, err error) {

	// execute the query
	profileQuery := "SELECT profile.uuid, profile.email, profile.firstname, profile.lastname,  CONCAT(profile.firstname, ' ', profile.lastname) AS fullname, profile.mobile, profile.birthdate, profile.gender, profile.updated, profile.created, CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', address.uuid, 'primary', profile_has_address.primary, 'street', address.street, 'additional', address.additional, 'zip', address.zip, 'city', address.city, 'country', address.country, 'google_id', address.google_id, 'updated', address.updated, 'created', address.created)) " +
		", ']') " +
		"FROM profile " +
		"LEFT JOIN profile_has_address ON profile.id = profile_has_address.profile_id " +
		"LEFT JOIN address ON address.id = profile_has_address.address_id " +
		"WHERE profile.email LIKE ? " +
		"GROUP BY profile.id " +
		sort + " " +
		"LIMIT ?, ?"

	rows, err := utils.DB.Query(profileQuery, filter.Email, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	var addressByte []byte

	// convert each row
	for rows.Next() {

		// create models
		profile := new(models.ProfileDefault)
		address := new(models.AddressList)

		// scan row and fill Profile
		err = rows.Scan(&profile.Uuid, &profile.Email, &profile.FirstName, &profile.LastName, &profile.FullName, &profile.Mobile, &profile.Birthdate, &profile.Gender, &profile.Updated, &profile.Created, &addressByte)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add addresses
		err = json.Unmarshal(addressByte, &address)

		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// add addresses to user
		if (*address)[0].Uuid != "" {
			profile.Addresses = *address.Distinct()
		} else {
			profile.Addresses = nil
		}

		Profiles = append(Profiles, *profile)
	}
	return Profiles, err
}

/**
 * select list of Profile with addresses and avatar and primary crew name
 */
func GetProfileMinList(page *models.Page, sort string, filter *models.FilterProfile) (Profiles []models.ProfileMin, err error) {

	// execute the query
	profileQuery := "SELECT profile.uuid, profile.email, CONCAT(profile.firstname, ' ', profile.lastname) AS fullname, profile.mobile, profile.birthdate, profile.gender, profile.created, " +
		"IFNULL(avatar.url, ''), IFNULL(avatar.type, ''), " +
		"IFNULL(crew.name, '') " +
		"FROM profile " +
		"LEFT JOIN avatar ON avatar.profile_id = profile.id " +
		"LEFT JOIN profile_has_crew ON profile_has_crew.profile_id = profile.id " +
		"LEFT JOIN nvm_state ON nvm_state.profile_has_crew_id = profile_has_crew.id AND nvm_state.primary_crew = 1 " +
		"LEFT JOIN crew ON profile_has_crew.crew_id = crew.id  " +
		"WHERE profile.email LIKE ? " +
		"GROUP BY profile.id " +
		sort + " " +
		"LIMIT ?, ?"

	rows, err := utils.DB.Query(profileQuery, filter.Email, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	// convert each row
	for rows.Next() {

		//create models
		profile := new(models.ProfileMin)
		avatar := new(models.Avatar)

		// scan row and fill Profile
		err = rows.Scan(&profile.Uuid, &profile.Email, &profile.FullName, &profile.Mobile, &profile.Birthdate, &profile.Gender, &profile.Created, &avatar.Url, &avatar.Type, &profile.PrimaryCrew)
		if err != nil {
			log.Print("Database Error: ", err)
			return nil, err
		}

		// Add avatar
		profile.Avatar = *avatar

		Profiles = append(Profiles, *profile)
	}
	return Profiles, err
}

/**
 * select list of Profile with addresses, avatars, crews and roles
 */
func GetProfileExtendedList(page *models.Page, sort string, filter *models.FilterProfile) (Profiles []models.ProfileExtended, err error) {

	// execute the query
	profileQuery := "SELECT profile.uuid, profile.email, profile.firstname, profile.lastname,  CONCAT(profile.firstname, ' ', profile.lastname) AS fullname, profile.mobile, profile.birthdate, profile.gender, newsletter.setting, profile.updated, profile.created, " +
		"IFNULL(avatar.url, ''), IFNULL(avatar.type, ''), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', address.uuid, 'primary', profile_has_address.primary, 'street', address.street, 'additional', address.additional, 'zip', address.zip, 'city', address.city, 'country', address.country, 'google_id', address.google_id, 'updated', address.updated, 'created', address.created)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', crew.uuid, 'primary', nvm_state.primary_crew, 'name', crew.name, 'email', crew.email, 'abbreviation', crew.abbreviation)) " +
		", ']'), " +
		"CONCAT('[', " +
		"GROUP_CONCAT(JSON_OBJECT('uuid', profile_has_crew_has_crew_role.uuid,'crew_uuid', crew.uuid, 'crew_name', crew.name, 'role', crew_role.name, 'description', crew_role.description)) " +
		", ']') " +
		"FROM profile " +
		"LEFT JOIN avatar ON avatar.profile_id = profile.id " +
		"LEFT JOIN profile_has_crew ON profile_has_crew.profile_id = profile.id " +
		"LEFT JOIN profile_has_crew_has_crew_role ON profile_has_crew_has_crew_role.profile_has_crew_id = profile_has_crew.id " +
		"LEFT JOIN crew_role ON crew_role.id = profile_has_crew_has_crew_role.crew_role_id " +
		"LEFT JOIN nvm_state ON nvm_state.profile_has_crew_id = profile_has_crew.id " +
		"LEFT JOIN crew ON profile_has_crew.crew_id = crew.id  " +
		"LEFT JOIN newsletter ON newsletter.profile_id = profile.id  " +
		"LEFT JOIN profile_has_address ON profile.id = profile_has_address.profile_id " +
		"LEFT JOIN address ON address.id = profile_has_address.address_id " +
		"WHERE profile.email LIKE ? " +
		"GROUP BY profile.id " +
		sort + " " +
		"LIMIT ?, ?"

	log.Print(profileQuery)

	rows, err := utils.DB.Query(profileQuery, filter.Email, page.Offset, page.Count)
	if err != nil {
		log.Print("Database Error", err)
		return nil, err
	}

	// byte arrays for group concatted jsons from sql
	var addressByte []byte
	var crewsByte []byte
	var rolesByte []byte

	// convert each row
	for rows.Next() {

		//create models
		profile := new(models.ProfileExtended)
		avatar := new(models.Avatar)
		crews := new(models.CrewList)
		address := new(models.AddressList)
		roles := new(models.ProfileRoleList)

		//scan row and fill Profile
		err = rows.Scan(&profile.Uuid, &profile.Email, &profile.FirstName, &profile.LastName, &profile.FullName, &profile.Mobile, &profile.Birthdate, &profile.Gender, &profile.Newsletter, &profile.Updated, &profile.Created, &avatar.Url, &avatar.Type, &addressByte, &crewsByte, &rolesByte)
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
		// add addresses to user
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
