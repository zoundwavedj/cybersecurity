package middlewares

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zoundwavedj/cybersecurity/configs"
	"github.com/zoundwavedj/cybersecurity/handlers"
)

// JwtMiddleware to handle authenticated routes
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			errorMsg            = "You're not authorized to access this resource"
			authorized          = false
			authorizationString = strings.TrimSpace(r.Header.Get("Authorization"))
		)

		if authorizationString != "" && strings.HasPrefix(authorizationString, "Bearer") {
			splits := strings.Split(authorizationString, " ")

			if len(splits) == 2 {
				if err := configs.ValidateAccessToken(splits[1]); err != nil {
					if err == configs.ErrTokenExpired || err == configs.ErrTokenMissing {
						errorMsg = err.Error()
					} else {
						log.Err(err).Msg("")
					}
				} else {
					authorized = true
				}
			}
		}

		if authorized {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			handlers.HandleError(w, errorMsg, http.StatusUnauthorized)
		}
	})
}
