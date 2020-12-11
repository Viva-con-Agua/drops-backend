package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

//	"database/sql"

type (
	//ProfileQuery represents query string for profile filter operations.
	ProfileQuery struct {
		UserID []string `query:"user_id"`
	}
)

//Filter creates mongodb bson filter from q.
func (q *ProfileQuery) Filter() bson.M {
	/*//inner query
	var query []bson.M

	//convert user_id list
	if q.UserID != nil {
		query = append(query, bson.M{"user_id": bson.M{"$or": q.UserID}})
	}

	filter := bson.M{"$filter": func() bson.M {
		if query != nil {
			if len(query) > 0 {
				return bson.M{"$and": query}
			}
		}
		return bson.M{}
	}()}
	log.Print(filter)*/
	query := bson.M{}
	if q.UserID != nil {
		query["user_id"] = bson.M{"$in": q.UserID}
	}
	return query
}
