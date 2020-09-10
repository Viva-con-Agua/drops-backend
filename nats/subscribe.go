package nats

import (
	"drops-backend/database"
	"drops-backend/models"
	"log"
)

func SubscribeAddModel() {
	_, err := Nats.Subscribe("drops.model.add", func(a *models.ModelCreate) {
		err := database.ModelInsert(a)
		if err != nil {
			log.Print("Database Error: ", err)
		}
	})
	if err != nil {
		log.Print("Nats Error: ", err)
	}
}

func SubscribeAccessAdd() {
	_, err := Nats.Subscribe("auth.access.add", func(a *models.AccessCreate) {
		err := database.AccessInsert(a)
		if err != nil {
			log.Print("Database Error: ", err)
		}
	})
	if err != nil {
		log.Print("Nats Error: ", err)
	}

}

func SubscribeAccessDelete() {
	_, err := Nats.Subscribe("auth.access.delete", func(a *models.DeleteBody) {
		err := database.AccessDelete(a)
		if err != nil {
			log.Print("Database Error: ", err)
		}
	})
	if err != nil {
		log.Print("Nats Error: ", err)
	}

}
