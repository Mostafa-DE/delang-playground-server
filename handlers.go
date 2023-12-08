package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func jsonHandler(resW http.ResponseWriter, req *http.Request, res map[string]string, filename string) {
	resW.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resW).Encode(res)

	if filename != "" {
		os.Remove(filename)
	}
}

func methodNotAllowedHandler(resW http.ResponseWriter, req *http.Request) {
	resW.WriteHeader(http.StatusMethodNotAllowed)
	res := map[string]string{
		"error": "Method not allowed",
	}

	jsonHandler(resW, req, res, "")
}
