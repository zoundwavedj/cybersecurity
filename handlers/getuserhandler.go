package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/database"
	"github.com/zoundwavedj/cybersecurity/utils"
)

type getUserResp struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Dob   string `json:"dob,omitempty"`
	Email string `json:"email,omitempty"`
	Ssn   string `json:"ssn,omitempty"`
}

// GetUserHandler function to retrieve a single user given ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		HandleError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := keys[0]

	rows, err := database.Db.Query("SELECT * FROM user WHERE id=?", id)
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		HandleError(w, "User not found", http.StatusNotFound)
		return
	}

	var (
		name  string
		dob   string
		email string
		ssn   string
	)

	if err = rows.Scan(&id, &name, &dob, &email, &ssn); err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	decryptedSsn, err := utils.Decrypt(ssn)
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	json.NewEncoder(w).Encode(&getUserResp{
		ID:    id,
		Name:  name,
		Dob:   dob,
		Email: email,
		Ssn:   decryptedSsn,
	})
}
