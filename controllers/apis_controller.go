package controllers

import (
	"net/http"
	"os"

	"github.com/gosu-team/cfapp-api/lib"
)

// APIInfo ...
type APIInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// GetAPIInfo ...
func GetAPIInfo(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	res.SendOK(APIInfo{
		Name:    os.Getenv("API_TITLE"),
		Version: os.Getenv("API_VERSION"),
	})
}
