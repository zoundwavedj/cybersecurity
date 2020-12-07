package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/configs"
	"github.com/zoundwavedj/cybersecurity/database"
)

type refreshTokenReq struct {
	Token string `json:"token,omitempty"`
}

type refreshTokenResp struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

var (
	accessID      string
	accessExpiry  int64
	accessToken   string
	refreshID     string
	refreshExpiry int64
	refreshToken  string
)

// RefreshTokenHandler func to recreate an expired token
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req refreshTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, oldRefreshID, err := configs.ValidateRefreshToken(req.Token)
	if err != nil {
		if err == configs.ErrTokenExpired || err == configs.ErrTokenMissing {
			HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	accessID, accessExpiry, accessToken, err = configs.GenerateAccessToken()
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	refreshID, refreshExpiry, refreshToken, err = configs.GenerateRefreshToken()
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	statement, err := database.Db.Prepare("INSERT INTO authentication (userId, token, expiresAt) VALUES (?, ?, ?)")
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, accessID, accessExpiry); err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	if _, err = statement.Exec(userID, refreshID, refreshExpiry); err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	statement.Close()
	statement, err = database.Db.Prepare("DELETE FROM authentication WHERE token=?")
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer statement.Close()

	if _, err = statement.Exec(oldRefreshID); err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	resp := &refreshTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	json.NewEncoder(w).Encode(resp)
}
