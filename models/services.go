package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	ServiceCreate struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	Service struct {
		Uuid        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Created     int64  `json:"created"`
		Updated     int64  `json:"updated"`
	}
)

func (service_create *ServiceCreate) Service() *Service {
	service := new(Service)
	service.Uuid = uuid.New().String()
	service.Name = service_create.Name
	service.Description = service_create.Description
	service.Created = time.Now().Unix()
	service.Updated = service.Created
	return service
}
