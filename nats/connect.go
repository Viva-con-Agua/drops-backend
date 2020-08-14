package nats

import (
	"log"

	"github.com/Viva-con-Agua/echo-pool/config"
	nats "github.com/nats-io/nats.go"
)

var Nats = new(nats.EncodedConn)

func Connect() {
	log.Print(config.Config.Nats.Url)
	nc, err := nats.Connect(config.Config.Nats.Url)
	if err != nil {
		log.Fatal("nats connection failed", err)
	}
	Nats, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal("nats encoded connection failed", err)
	}
	log.Print("nats successfully connected!")
}
