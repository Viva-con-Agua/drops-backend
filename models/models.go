package models

type (
	// Model can be every model
	Model struct {
		Uuid        string `json:"uuid" validate:"required"`
		Type        string `json:"type" validate:"required"`
		Name        string `json:"name" validate:"required"`
		ServiceName string `json:"service_name" validate:"required"`
		Owner       string `json:"owner" validate:"required"`
		Created     string `json:"created" validate:"required"`
	}
	ModelStub struct {
		Uuid        string `json:"uuid" validate:"required"`
		Type        string `json:"type" validate:"required"`
		Name        string `json:"name" validate:"required"`
		ServiceName string `json:"service_name" validate:"required"`
		Owner       string `json:"owner" validate:"required"`
	}
)
