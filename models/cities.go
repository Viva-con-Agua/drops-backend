package models

type (
	CityCreate struct {
		City     string `json:"city" validate;"required"`
		Country  string `json:"country" validate;"required"`
		GoogleId string `json:"google_id" validate;"required"`
	}
	City struct {
		Uuid     string `json:"uuid" validate:"required"`
		City     string `json:"city" validate:"required"`
		Country  string `json:"country" validate:"required"`
		GoogleId string `json:"google_id" validate:"required"`
	}
	CityList []City
)

func (list *CityList) Distinct() *CityList {
	r := make(CityList, 0, len(*list))
	m := make(map[City]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}
