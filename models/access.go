package models

type (
	Access struct {
		Uuid        string `json:"uuid" validate:"required"`
		AccessName  string `json:"access_name" validate:"required"`
		ServiceName string `json:"service_name" validate:"required"`
		ModelUuid   string `json:"model_uuid"`
		ModeName    string `json:"model_name"`
		ModelType   string `json:"model_type"`
		Created     int    `json:"created" validate:"required"`
	}
	AccessCreate struct {
		Assign    string `json:"assign"`
		Name      string `json:"name"`
		ModelUuid string `json:"model_uuid"`
	}
)
