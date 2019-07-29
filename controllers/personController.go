package controllers

import (
	"../models"
	u "../utils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func CreatePerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	person := &models.Person{}
	err := json.NewDecoder(request.Body).Decode(person)
	if err != nil {
		u.Respond(response, u.Message(false, "Error while decoding request body"))
		return
	}

	collection := models.GetDB().Database("go-rest-api-case").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.InsertOne(ctx, person)

	json.NewEncoder(response).Encode(person)
}

func GetPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := models.GetDB().Database("go-rest-api-case").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	person := &models.Person{}
	err := collection.FindOne(ctx, &models.Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(person)
}

func GetPeople(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []models.Person
	collection := models.GetDB().Database("go-rest-api-case").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person models.Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func DeletePerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := models.GetDB().Database("go-rest-api-case").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	person := &models.Person{}
	collection.DeleteOne(ctx, &models.Person{ID: id})

	json.NewEncoder(response).Encode(person)
}