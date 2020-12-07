package handlers

import "net/http"

// NotFoundHandler function to handle unregistered routes
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Oops, the resource you're looking for doesn't exist m8", http.StatusNotFound)
}
