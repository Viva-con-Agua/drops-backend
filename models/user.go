package models

type (
	//UserQuery represents query string for user database filter
	UserQuery struct {
		Offset      int      `query:"offset" default:"0"`
		Count       int      `query:"count" default:"40"`
		SortDir     string   `query:"sortdir"`
		SortBy      string   `query:"sortby"`
		Email       []string `query:"email" `
		UpdatedFrom int      `query:"updated_from"`
		UpdatedTo   int      `query:"updated_to"`
	}
)
