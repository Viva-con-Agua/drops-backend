package models

type (
	AddressCreate struct {
		Street     string `json:"street" validate:"required"`
		Additional string `json:"additional" validate;"required"`
		Zip        string `json:"zip" validate;"required"`
		City       string `json:"city" validate;"required"`
		Country    string `json:"country" validate;"required"`
		GoogleId   string `json:"google_id" validate;"required"`
		Updated    int    `json:"updated" validate:"required"`
		Created    int    `json:"created" validate:"required"`
	}
	AddressUpdate struct {
		Uuid       string `json:"publicId" validate;"required"`
		Street     string `json:"street" validate:"required"`
		Additional string `json:"additional" validate;"required"`
		Zip        string `json:"zip" validate;"required"`
		City       string `json:"city" validate;"required"`
		Country    string `json:"country" validate;"required"`
		Primary    bool   `json:"primary" validate;"required"`
		GoogleId   string `json:"google_id" validate;"required"`
	}

	Address struct {
		Uuid       string `json:"uuid" validate:"required"`
		Primary    int    `json:"primary" validate:"required"`
		Street     string `json:"street" validate:"required"`
		Additional string `json:"additional" validate:"required"`
		Zip        string `json:"zip" validate:"required"`
		City       string `json:"city" validate:"required"`
		Country    string `json:"country" validate:"required"`
		GoogleId   string `json:"google_id" validate:"required"`
		Updated    int    `json:"updated" validate:"required"`
		Created    int    `json:"created" validate:"required"`
	}
	AddressList []Address
)

func (list *AddressList) Distinct() *AddressList {
	r := make(AddressList, 0, len(*list))
	m := make(map[Address]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
}
