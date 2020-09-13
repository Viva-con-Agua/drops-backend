package models

import (
	"log"

	"github.com/Viva-con-Agua/echo-pool/api"
)

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
	temp := make(map[string][]string)
	log.Print(as)
	temp[(*as)[0].Model] = append(temp[(*as)[0].Model], (*as)[0].Right)
	resp[(*as)[0].Service] = temp
	for _, s := range (*as)[1:] {
		resp[s.Service][s.Model] = append(resp[s.Service][s.Model], s.Right)
	}
	return resp
}

func (access_db *AccessDB) Access() *api.Access {
	access := new(api.Access)
	access.AccessUuid = access_db.AccessUuid
	access.AccessName = access_db.AccessName
	access.ModelUuid = access_db.ModelUuid
	access.ModelName = access_db.ModelName
	access.ModelType = access_db.ModelType
	return access
}

func (list *AccessDBList) Distinct() *AccessDBList {
	r := make(AccessDBList, 0, len(*list))
	m := make(map[AccessDB]bool)
	for _, val := range *list {
		if _, ok := m[val]; !ok {
			m[val] = true
			r = append(r, val)
		}
	}
	return &r
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
