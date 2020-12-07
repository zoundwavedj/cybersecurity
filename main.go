package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/database"
	"github.com/zoundwavedj/cybersecurity/handlers"
	"github.com/zoundwavedj/cybersecurity/middlewares"
)

func main() {
	// Log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting up...")

	database.Cleanup()
	database.Setup()
	defer database.Db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/superuser", handlers.CreateSuperUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.UserLoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", handlers.UserLogoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods(http.MethodPost)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	ar := r.NewRoute().Subrouter()
	ar.Use(middlewares.JwtMiddleware)
	ar.HandleFunc("/user", handlers.CreateUserHandler).Methods(http.MethodPost)
	ar.HandleFunc("/user", handlers.GetUserHandler).Methods(http.MethodGet).Queries("id", "")
	ar.HandleFunc("/users", handlers.ListUsersHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("")
		}
	}()

	log.Info().Msg("Cybersec Test application is ready :D")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	log.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 secs to wait for all background work to finish, if any
	defer cancel()
	srv.Shutdown(ctx)

	log.Info().Msg("Bye bye :D")
	os.Exit(0)
}
