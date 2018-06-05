package endpoint

import (
	"net/http"
	"encoding/json"
	"connector"
)

func GetOpenGame(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(connector.GetOpenGame(-1))
}
