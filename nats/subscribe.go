package nats

/*import (
	"a-backend/database"
	"auth-backend/models"
	"log"
)

func SubscribeAccessAdd() {
	_, err := Nats.Subscribe("auth.access.add", func(a *models.AccessUserCreate) {
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

}*/
