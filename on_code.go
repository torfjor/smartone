package smartone

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// OnCodeRequest represents an incoming code request from Tidomat smartONE
type OnCodeRequest struct {
	// Code typed at reader
	Code string `json:"Code"`
	// ID of door that the Code was used at
	DoorID       int    `json:"DoorId"`
	SLXID        string `json:"SLXId"`
	SerialNumber string `json:"SERNO"`
}

// OnCodeResponse represents a response to Tidomat smartONE in reply to a
// OnCodeRequest.
// A response with Result set to true instructs the Tidomat smartONE that this
// was a successful request.
type OnCodeResponse struct {
	Result bool `json:"result"`
	// Optional. Message to log in case of Result is false.
	ResultMessage string `json:"resultMessage,omitempty"`
}

type OnCodeFunc func(ctx context.Context, req OnCodeRequest) (OnCodeResponse, error)

func OnCodeHandler(fn OnCodeFunc) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Add("Allow", http.MethodPost)
			writeErr(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			writeErr(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var req OnCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		res, err := fn(r.Context(), req)
		if err != nil {
			writeErr(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(res)
		if err != nil {
			writeErr(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
	}

	return handler
}
