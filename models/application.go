package models

import (
	"time"

	"github.com/Viva-con-Agua/vcago/vmod"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	//ApplicationCreate represents a struct for application create
	ApplicationCreate struct {
		Name      string        `bson:"name" json:"name" validate:"required"`
		CreatorID string        `bson:"creator_id" json:"creator_id" validate:"required"`
		Modified  vmod.Modified `bson:"modified" json:"modified"`
	}
	//Application represents a frontend application
	Application struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"app_id" validate:"required"`
		Name      string             `bson:"name" json:"name" validate:"required"`
		CreatorID string             `bson:"creator_id" json:"creator_id" validate:"required"`
		Modified  vmod.Modified      `bson:"modified" json:"modified"`
	}
)

func (a *ApplicationCreate) Application() *Application {
	return &Application{
		Name:      a.Name,
		CreatorID: a.CreatorID,
		Modified:  *vmod.NewModified(time.Now().Unix()),
	}
}
