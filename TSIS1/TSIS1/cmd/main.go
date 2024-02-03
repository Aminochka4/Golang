package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "TSIS1/pkg"
)

func main(){
    log.Println("Starting API server")
    router := mux.NewRouter()
    router.HandleFunc("/health-check", pkg.HealthCheck).Methods("GET")
    router.HandleFunc("/sport-teams", pkg.SportTeams).Methods("GET")
    router.HandleFunc("/sport-teams/{id:[0-6]+}", pkg.SportTeamsByID).Methods("GET")

    http.ListenAndServe(":8080", router)
}