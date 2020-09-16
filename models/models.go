package models

type (
	// Model can be every model
	Model struct {
		Uuid    string `json:"uuid" validate:"required"`
		Name    string `json:"name" validate:"required"`
		Service string `json:"service" validate:"required"`
		Creator string `json:"owner" validate:"required"`
		Created int64  `json:"created" validate:"required"`
	}
	ModelCreate struct {
		Uuid    string `json:"uuid" validate:"required"`
		Name    string `json:"name" validate:"required"`
		Service string `json:"service" validate:"required"`
		Creator string `json:"owner" validate:"required"`
	}
)

func (m_create *ModelCreate) Model(created int64) *Model {
	m := new(Model)
	m.Uuid = m_create.Uuid
	m.Name = m_create.Name
	m.Service = m_create.Service
	m.Creator = m_create.Creator
	m.Created = created
	return m
}
