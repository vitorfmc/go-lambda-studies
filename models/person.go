package models

import (
	u "../utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func (person *Person) Validate() (map[string] interface{}, bool) {

	if person.Firstname == "" {
		return u.Message(false, "Firstname should be on the payload"), false
	}

	return u.Message(true, "success"), true
}