package database

import (
	"os"

	"github.com/rs/zerolog/log"
)

// Cleanup function to remove previous database file
func Cleanup() {
	if err := os.Remove("./local.db"); err != nil {
		log.Fatal().Err(err)
	}
}
