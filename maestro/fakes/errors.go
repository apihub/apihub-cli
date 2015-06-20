package fakes

import (
	"encoding/json"
	"net/http"
)

func (fake *BackstageServer) Error(w http.ResponseWriter, statusCode int, i interface{}) {
	j, _ := json.Marshal(i)
	w.WriteHeader(statusCode)
	w.Write(j)
}
