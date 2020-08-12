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
	Dict struct {
		Key   string `json:"key" validate:"required"`
		Value string `json:"value" validate:"required"`
	}
	DictList []Dict
)

func (list *DictList) Distinct() *DictList {
	r := make(DictList, 0, len(*list))
	m := make(map[Dict]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}
