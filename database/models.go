package database

import (
	"drops-backend/models"
	"drops-backend/utils"
	"log"
	"time"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/google/uuid"
)

func ModelCreate(m_create *api.ModelCreate) (m *api.Model, ap_err *api.ApiError) {
	created := time.Now().Unix()
	tx, err := utils.DB.Begin()
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	log.Print(m_create)
	// select drops_user database id as u_id
	query := "SELECT id FROM drops_user WHERE uuid = ? LIMIT 1"
	var u_id int64
	err = tx.QueryRow(query, m_create.Creator).Scan(&u_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	// insert model
	query = "INSERT INTO model (uuid, name, service, created, updated) " +
		"VALUES(?, ?, ?, ?, ?)"
	res, err := tx.Exec(query, m_create.Uuid, m_create.Name, m_create.Service, created, created)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	m_id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	log.Print(m_id)
	query = "INSERT INTO access " +
		"(uuid, service, updated, created, drops_user_id, model_id)" +
		" VALUES(?, ?, ?, ?, ?, ?)"
	res, err = tx.Exec(query, uuid.New(), m_create.Service, created, created, u_id, m_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	a_id, err := res.LastInsertId()
	log.Print(a_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}

	query = "INSERT INTO `access_right` ( `right`, created, `access_id`) " +
		"VALUES('created', ?, ?)"
	_, err = tx.Exec(query, created, a_id)
	if err != nil {
		tx.Rollback()
		return nil, api.GetError(err)
	}
	m = m_create.Model(created)
	return m, api.GetError(tx.Commit())
}

func ModelDelete(deleteBody *models.DeleteBody) (err error) {
	tx, err := utils.DB.Begin()
	if err != nil {
		log.Print("Database Error: ", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM model WHERE uuid = ?", deleteBody.Uuid)
	if err != nil {
		tx.Rollback()
		log.Print("Database Error: ", err)
		return err
	}
	return tx.Commit()
}
