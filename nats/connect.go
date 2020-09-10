package nats

import (
	"log"
	"os"

	nats "github.com/nats-io/nats.go"
)

var Nats = new(nats.EncodedConn)

func Connect() {
	nc, err := nats.Connect("nats://" + os.Getenv("NATS_HOST") + ":" + os.Getenv("NATS_PORT"))
	if err != nil {
		log.Fatal("nats connection failed", err)
	}
	Nats, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal("nats encoded connection failed", err)
	}
	log.Print("nats successfully connected!")
}
