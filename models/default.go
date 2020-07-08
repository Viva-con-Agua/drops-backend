package models

type (
	Page struct {
		Offset int
		Count  int
	}
	DeleteBody struct {
		Uuid string `json:"uuid" validate:"required"`
	}
	AssignBody struct {
		Assign string `json:"assign" validate:"required"`
		To     string `json:"to" validate:"required"`
	}
)
