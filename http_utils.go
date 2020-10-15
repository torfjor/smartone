package smartone

import "net/http"

func writeErr(w http.ResponseWriter, msg string, status int) error {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
