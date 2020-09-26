package models

type (
	// For handling database

	AccessDB struct {
		AccessUuid  string `json:"access_uuid" validate:"required"`
		AccessName  string `json:"access_name" validate:"required"`
		ServiceName string `json:"service_name" validate:"required"`
		ModelUuid   string `json:"model_uuid"`
		ModelName   string `json:"model_name"`
		ModelType   string `json:"model_type"`
		Created     int    `json:"created" validate:"required"`
	}
	AccessDBList []AccessDB

	AccessCreate struct {
		Assign    string `json:"assign"`
		Name      string `json:"name"`
		ModelUuid string `json:"model_uuid"`
	}
	AccessDefault struct {
		ServiceName string
		UserUuid    string
		AccessType  string
	}
	AccessSession struct {
		Service string   `json:"service"`
		Rights  []string `json:"rights"`
		Model   string   `json:"model"`
	}
	AccessSessionDB struct {
		Service string `json:"service"`
		Model   string `json:"model"`
		Right   string `json:"right"`
	}
	AccessSessionDBList []AccessSessionDB
	AccessString        map[string]map[string]string
)

func (as *AccessSessionDBList) List() map[string]map[string][]string {
	resp := make(map[string]map[string][]string)
	for _, s := range *as {
		if len(resp[s.Service]) == 0 {
			sign := make(map[string][]string)
			sign[s.Model] = append(sign[s.Model], s.Right)
			resp[s.Service] = sign
		} else {
			resp[s.Service][s.Model] = append(resp[s.Service][s.Model], s.Right)
		}
	}
	return resp
}

/*
func (list *AccessDBList) AccessList() *api.AccessList {
	d_list := list.Distinct()
	access_list := make(api.AccessList)
	for _, val := range *d_list {
		access_list[val.ServiceName] = append(access_list[val.ServiceName], *val.Access())
	}
	return &access_list

}*/
