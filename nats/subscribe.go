package nats

/*
func SubscribeAddModel() {
	_, err := Nats.Subscribe("drops.model.add", func(a *api.ModelCreate) {
		_, api_err := database.ModelCreate(a)
		if api_err != nil {
			log.Print(api_err)
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

}*/
