package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/database"
	"github.com/zoundwavedj/cybersecurity/utils"
)

type createUserReq struct {
	Name  string `json:"name,omitempty"`
	Dob   string `json:"dob,omitempty"`
	Email string `json:"email,omitempty"`
	Ssn   string `json:"ssn,omitempty"`
}

type createUserResp struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Dob   string `json:"dob,omitempty"`
	Email string `json:"email,omitempty"`
	Ssn   string `json:"ssn,omitempty"`
}

// CreateUserHandler function to create users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req createUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	statement, err := database.Db.Prepare("INSERT INTO user (id, name, dob, email, ssn) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer statement.Close()

	id := uuid.New().String()
	encryptedSsn, err := utils.Encrypt(req.Ssn)
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	if _, err := statement.Exec(id, req.Name, req.Dob, req.Email, encryptedSsn); err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	json.NewEncoder(w).Encode(&createUserResp{
		ID:    id,
		Name:  req.Name,
		Dob:   req.Dob,
		Email: req.Email,
		Ssn:   encryptedSsn,
	})
}
