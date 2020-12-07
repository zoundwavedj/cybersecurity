package configs

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zoundwave/cybersecurity/database"
)

// Claims type
type Claims struct {
	AccessID  string `json:"accessId,omitempty"`
	RefreshID string `json:"refreshId,omitempty"`
	jwt.StandardClaims
}

// SecretType enum type
type SecretType string

const (
	// ACCESS enum
	ACCESS SecretType = "ACCESS"
	// REFRESH enum
	REFRESH SecretType = "REFRESH"
)

var (
	// ErrTokenMissing error
	ErrTokenMissing = errors.New("Token missing")
	// ErrTokenExpired error
	ErrTokenExpired = errors.New("Token expired")
	accessSecret    = []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret   = []byte(os.Getenv("REFRESH_SECRET"))
)

// GenerateAccessToken function to build access token given ID
func GenerateAccessToken() (string, int64, string, error) {
	uuid := uuid.New().String()
	expiresAt := time.Now().Add(time.Second * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		AccessID: uuid,
	})

	signedToken, err := token.SignedString(accessSecret)

	return uuid, expiresAt, signedToken, err
}

// GenerateRefreshToken function to build refresh token
func GenerateRefreshToken() (string, int64, string, error) {
	uuid := uuid.New().String()
	expiresAt := time.Now().Add(time.Minute * 5).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		RefreshID: uuid,
	})

	signedToken, err := token.SignedString(refreshSecret)

	return uuid, expiresAt, signedToken, err
}

// ValidateAccessToken function to verify JWT token
func ValidateAccessToken(tokenString string) error {
	token, err := ParseToken(tokenString, ACCESS)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return err
	}

	accessID, ok := claims["accessId"].(string)
	if !ok {
		return err
	}

	rows, err := database.Db.Query("SELECT userId, expiresAt FROM authentication WHERE token=? LIMIT 1", accessID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var (
			userID    string
			expiresAt int64
		)

		if err = rows.Scan(&userID, &expiresAt); err != nil {
			return err
		}
		rows.Close()

		if time.Now().After(time.Unix(expiresAt, 0)) {
			statement, err := database.Db.Prepare("DELETE FROM authentication WHERE token=?")
			if err != nil {
				return err
			}
			defer statement.Close()

			if _, err = statement.Exec(accessID); err != nil {
				return err
			}

			return ErrTokenExpired
		}
	} else {
		return ErrTokenMissing
	}

	return nil
}

// ValidateRefreshToken function to verify JWT token
func ValidateRefreshToken(tokenString string) (string, string, error) {
	token, err := ParseToken(tokenString, REFRESH)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", err
	}

	refreshID, ok := claims["refreshId"].(string)
	if !ok {
		return "", "", err
	}

	rows, err := database.Db.Query("SELECT userId, expiresAt FROM authentication WHERE token=? LIMIT 1", refreshID)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var userID string

	if rows.Next() {
		var expiresAt int64

		if err = rows.Scan(&userID, &expiresAt); err != nil {
			return "", "", err
		}
		rows.Close()

		if time.Now().After(time.Unix(expiresAt, 0)) {
			statement, err := database.Db.Prepare("DELETE FROM authentication WHERE token=?")
			if err != nil {
				return "", "", err
			}
			defer statement.Close()

			if _, err = statement.Exec(refreshID); err != nil {
				return "", "", err
			}

			return "", "", ErrTokenExpired
		}
	} else {
		return "", "", ErrTokenMissing
	}

	return userID, refreshID, nil
}

// ParseToken function to parse a token string
func ParseToken(token string, secretType SecretType) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		if secretType == ACCESS {
			return accessSecret, nil
		}

		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
