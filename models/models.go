package models

type (
	// Model can be every model
	Model struct {
		Uuid 		string `json:"uuid" validate:"required"`
		Type    string `json:"type" validate:"required"`
		Updated string `json:"updated" validate:"required"`
		Created string `json:"created" validate:"required"`

	}
)
