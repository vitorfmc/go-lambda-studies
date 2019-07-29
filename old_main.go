package main

import (
	"./controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/person", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/v1/person", controllers.GetPeople).Methods("GET")
	router.HandleFunc("/v1/person/{id}", controllers.GetPerson).Methods("GET")
	router.HandleFunc("/v1/person/{id}", controllers.DeletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}