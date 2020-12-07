package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/configs"
	"github.com/zoundwavedj/cybersecurity/database"
)

type userLogoutReq struct {
	AccessToken  string
	RefreshToken string
}

type userLogoutResp struct {
	Success bool `json:"success"`
}

// UserLogoutHandler function to handle user logouts
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req userLogoutReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var (
		claims    jwt.MapClaims
		accessID  string
		refreshID string
	)

	accessToken, err := configs.ParseToken(req.AccessToken, configs.ACCESS)
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	claims, _ = accessToken.Claims.(jwt.MapClaims)
	accessID, _ = claims["accessId"].(string)

	refreshToken, err := configs.ParseToken(req.RefreshToken, configs.REFRESH)
	if err != nil {
		log.Err(err).Msg("")
		HandleError500(w)
		return
	}

	claims, _ = refreshToken.Claims.(jwt.MapClaims)
	refreshID, _ = claims["refreshId"].(string)

	if accessID != "" && refreshID != "" {
		statement, err := database.Db.Prepare("DELETE FROM authentication WHERE token=?")
		if err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}
		defer statement.Close()

		if _, err = statement.Exec(accessID); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}

		if _, err = statement.Exec(refreshID); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}
	}

	json.NewEncoder(w).Encode(&userLogoutResp{
		Success: true,
	})
}
