package controllers

import (
	"encoding/json"
	"net/http"
)

// EncodeResponse turns the response into json and writes it back to the ResponseWriter
func EncodeResponse(w http.ResponseWriter, i interface{}) error {
	err := json.NewEncoder(w).Encode(&i)
	if err != nil {
		return err
	}
	return nil
}

// DecodeRequest decodes the request
func DecodeRequest(r *http.Request, i interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&i)
}
