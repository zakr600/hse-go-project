package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var (
	ContentTypeHeader = "Content-Type"
	JsonContent       = "application/json"
)

// TryParseJsonQuery deserialize json object in request
func TryParseJsonQuery(w http.ResponseWriter, r *http.Request, v any) bool {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		var output struct {
			ErrorMessage string `json:"errorMessage"`
		}
		output.ErrorMessage = err.Error()
		SendJson(w, r, output, http.StatusBadRequest)
		return false
	}
	return true
}

// SendJson sends http responce contains serialized object in json format
func SendJson(w http.ResponseWriter, r *http.Request, v any, httpCode int) {
	responseBody, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set(ContentTypeHeader, JsonContent)
	w.WriteHeader(httpCode)
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println(err)
	}

}
