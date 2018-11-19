package lib

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Request ...
type Request struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// GetJSONBody ...
func (r *Request) GetJSONBody(model interface{}) {
	decoder := json.NewDecoder(r.Request.Body)
	decoder.Decode(&model)
}

// GetVarID ...
func (r *Request) GetVarID() (int, error) {
	vars := mux.Vars(r.Request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.Info("Could not convert parameter id to int")
		http.Error(r.ResponseWriter, "Invalid id parameter", http.StatusBadRequest)
		return 0, err
	}
	return id, nil
}

// GetCommentPath ...
func (r *Request) GetCommentPath() (string, string, error) {
	vars := mux.Vars(r.Request)
	domain := vars["domain"]
	path := vars["path"]

	if domain == "" || path == "" {
		logrus.Info("Parameter is null, cannot convert")
		http.Error(r.ResponseWriter, "Invalid domain / path parameter", http.StatusBadRequest)
		return "", "", nil
	}
	return domain, path, nil
}