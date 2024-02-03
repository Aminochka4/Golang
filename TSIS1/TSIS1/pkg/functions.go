package pkg

import (
	"log"
	"net/http"
	"encoding/json"
	"TSIS1/api"
	"github.com/gorilla/mux"
)

func SportTeams(w http.ResponseWriter, r *http.Request) {
	log.Println("teams checking")
	sportTeams := api.Teams

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(sportTeams)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
    	return
	}
	w.Write(response)
}


func SportTeamsByID(w http.ResponseWriter, r *http.Request) {
	log.Println("team checking")

	vars := mux.Vars(r)
	teamID := vars["id"]

	sportTeam := api.GetTeam(teamID)
	w.Header().Set("Content-Type", "application/json")

	if sportTeam == nil {
		http.Error(w, "Such id of team does not exist", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(sportTeam)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
