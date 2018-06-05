package endpoint

import (
	"net/http"
	"connector"
	"github.com/gorilla/mux"
	"encoding/json"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if player := connector.GetPlayerByName(params["name"]); player.PlayerStatus != connector.OFFLINE {
		http.Error(w, "User already logged in", 403)
	} else {
		json.NewEncoder(w).Encode(connector.AddPlayer(params["name"]))
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	connector.LogoutPlayer(params["name"])
}
