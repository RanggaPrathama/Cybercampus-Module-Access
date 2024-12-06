package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModule struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	IDUser     primitive.ObjectID `json:"id_user" bson:"id_user"`
	JenisUser  string 			`json:"jenis_user" bson:"jenis_user"`
	MODULES    []primitive.ObjectID   `json:"modules" bson:"modules"`
	CREATED_AT time.Time             `json:"created_at" bson:"created_at"`
	UPDATED_AT time.Time `json:"updated_at" bson:"updated_at"`
}

type UserModuleResponse struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	
	Created_at time.Time `json:"created_at" bson:"created_at"`
	Updated_at time.Time `json:"updated_at" bson:"updated_at"`
}