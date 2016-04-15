package buddy

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(code int, w http.ResponseWriter, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	return err
}
