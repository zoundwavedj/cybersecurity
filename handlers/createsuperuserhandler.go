package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/configs"
	"github.com/zoundwavedj/cybersecurity/database"
	"golang.org/x/crypto/argon2"
)

type createSuperUserResp struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// CreateSuperUserHandler function to create a superuser
func CreateSuperUserHandler(w http.ResponseWriter, r *http.Request) {
	var username = os.Getenv("SUPERUSERNAME")
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.Db.Query("SELECT * FROM superuser")
	if err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer rows.Close()

	if rows.Next() {
		HandleError(w, "Superuser already exists", http.StatusBadRequest)
		return
	}

	hash, plaintext, err := generateHash()
	if err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}

	statement, err := database.Db.Prepare("INSERT INTO superuser (id, username, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}
	defer statement.Close()

	if _, err = statement.Exec(id.String(), username, string(hash)); err != nil {
		log.Error().Err(err).Msg("")
		HandleError500(w)
		return
	}

	resp := createSuperUserResp{
		Username: username,
		Password: plaintext,
	}
	json.NewEncoder(w).Encode(resp)
}

func generateHash() (string, string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}

	hashConfig := configs.DefaultHashConfig()

	salt := make([]byte, hashConfig.KeyLen)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	hash := argon2.IDKey([]byte(uuid.String()), salt, hashConfig.Time, hashConfig.Memory, hashConfig.Threads, hashConfig.KeyLen)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"

	return fmt.Sprintf(format, argon2.Version, hashConfig.Memory, hashConfig.Time, hashConfig.Threads, encodedSalt, encodedHash), uuid.String(), nil
}
