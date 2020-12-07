package handlers

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/configs"
	"github.com/zoundwavedj/cybersecurity/database"
	"golang.org/x/crypto/argon2"
)

type authenticateUserReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type authenticateUserResp struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// UserLoginHandler function handle user login and authentication
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req authenticateUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, "Make sure your request body is valid", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" && strings.TrimSpace(req.Password) == "" {
		HandleError(w, "Make sure your request body is valid", http.StatusBadRequest)
		return
	}

	rows, err := database.Db.Query("SELECT id, password FROM superuser WHERE username=? LIMIT 1", req.Username)
	if err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer rows.Close()

	var (
		accessID      string
		accessExpiry  int64
		accessToken   string
		refreshID     string
		refreshExpiry int64
		refreshToken  string
	)

	if rows.Next() {
		var id string
		var hashed string

		if err = rows.Scan(&id, &hashed); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}

		success, err := validateHash(req.Password, hashed)
		if err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}

		if !success {
			HandleError(w, "Invalid login", http.StatusBadRequest)
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

		if _, err = statement.Exec(id, accessID, accessExpiry); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}

		if _, err = statement.Exec(id, refreshID, refreshExpiry); err != nil {
			log.Err(err).Msg("")
			HandleError500(w)
			return
		}
	} else {
		HandleError(w, "Invalid login", http.StatusBadRequest)
		return
	}

	resp := authenticateUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	json.NewEncoder(w).Encode(resp)
}

func validateHash(plain string, hash string) (bool, error) {
	var hashConfig configs.HashConfig
	var version int

	split := strings.Split(hash, "$")

	_, err := fmt.Sscanf(split[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, errors.New("Hash version is incompatible")
	}

	_, err = fmt.Sscanf(split[3], "m=%d,t=%d,p=%d", &hashConfig.Memory, &hashConfig.Time, &hashConfig.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(split[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(split[5])
	if err != nil {
		return false, err
	}

	hashConfig.KeyLen = uint32(len(decodedHash))

	rehashed := argon2.IDKey([]byte(plain), salt, hashConfig.Time, hashConfig.Memory, hashConfig.Threads, hashConfig.KeyLen)

	return (subtle.ConstantTimeCompare(rehashed, decodedHash) == 1), nil
}
