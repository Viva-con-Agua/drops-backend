package controllers

import (
	"drops-backend/database"
	"drops-backend/utils"
	"log"
)

func AddEssential() (err error) {
	log.Print("Add Essential Models")
	err = database.AddEssentialModels()
	if err != nil {
		if err == utils.ErrorConflict {
			log.Print("Models already exsists")
			return nil
		}
		return err
	}
	return err
}
