package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/database"
)

type listUsersResp struct {
	Users []string `json:"users"`
}

// ListUsersHandler function to list all user IDs
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.Db.Query("SELECT id FROM user")
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer rows.Close()

	var resp listUsersResp
	resp.Users = []string{}

	for rows.Next() {
		var id string

		if err = rows.Scan(&id); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}

		resp.Users = append(resp.Users, id)
	}

	json.NewEncoder(w).Encode(resp)
}
