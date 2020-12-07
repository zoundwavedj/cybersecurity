package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResp struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

// HandleError function to write error resp as json
func HandleError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)

	resp := errorResp{
		Code:  code,
		Error: err,
	}

	json.NewEncoder(w).Encode(resp)
}

// HandleError500 function to handle default 500 code
func HandleError500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	resp := errorResp{
		Code:  http.StatusInternalServerError,
		Error: "Oops, something must've went wrong somewhere",
	}

	json.NewEncoder(w).Encode(resp)
}
