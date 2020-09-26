package models

type (
	AccessToken struct {
		Token   string `json:"token"`
		Tcase   string `json:"t_case"`
		expired int64  `json:"expired"`
		created int64  `json:"created"`
	}
)
