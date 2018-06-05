package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"endpoint"
)

func main() {
	router := mux.NewRouter()
	addHandlersToRouter(router);
	log.Fatal(http.ListenAndServe(":8080", router))
}

func addHandlersToRouter(router *mux.Router) {
	router.HandleFunc("/login/{name}", endpoint.HandleLogin).Methods("GET")
	router.HandleFunc("/logout/{name}", endpoint.HandleLogout).Methods("GET")
	router.HandleFunc("/game", endpoint.GetOpenGame).Methods("GET")
}