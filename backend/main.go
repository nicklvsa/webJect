package main

import (
	"net/http"
	"webject/router"
	"webject/shared"

	"github.com/gorilla/mux"
)

var (
	logger *shared.Logger
)

func buildHandlers() {
	r := mux.NewRouter()

	r.HandleFunc("/", router.BaseAPIHandler).Methods("GET", "POST")
	r.HandleFunc("/tweak/{action}", router.TweakAPIHandler).Methods("POST")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(router.NotFoundHandler).GetHandler()
	http.Handle("/", r)
}

func main() {
	buildHandlers()
	logger.Info("API started!")
	logger.Err(http.ListenAndServe(":8080", nil))
}
