package database

import (
	"drops-backend/models"
	"drops-backend/utils"

	"github.com/Viva-con-Agua/echo-pool/api"
)

func ProfileListInternal(filter string) (p_list []models.Profile, api_err *api.ApiError) {
	query := "SELECT " +
		"du.uuid, " +
		"p.uuid, p.first_name, p.last_name, CONCAT(p.first_name, ' ', p.last_name), " +
		"ifnull(p.display_name, 'none'), ifnull(p.gender, 'none'), p.updated, p.created, " +
		"ifnull(a.type, ''), ifnull(a.url,''), ifnull(a.updated, 0), ifnull(a.created, 0) " +
		"FROM drops_user AS du " +
		"LEFT JOIN profile AS p ON p.drops_user_id = du.id " +
		"LEFT JOIN avatar AS a ON a.profile_id = p.id " +
		filter

	rows, err := utils.DB.Query(query)
	if err != nil {
		return nil, api.GetError(err)
	}
	p := new(models.Profile)
	for rows.Next() {
		err = rows.Scan(
			&p.UserUuid,
			&p.Uuid,
			&p.FirstName,
			&p.LastName,
			&p.FullName,
			&p.DisplayName,
			&p.Gender,
			&p.Updated,
			&p.Created,
			&p.Avatar.Url,
			&p.Avatar.Type,
			&p.Avatar.Updated,
			&p.Avatar.Created,
		)
		if err != nil {
			return nil, api.GetError(err)
		}
		p_list = append(p_list, *p)
	}
	return p_list, api.GetError(err)

}
