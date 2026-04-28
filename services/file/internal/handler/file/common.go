package file

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type publicStatusError interface {
	error
	HTTPStatus() int
	PublicMessage() string
}

func writeDataPlaneError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	var dataErr publicStatusError
	if errors.As(err, &dataErr) {
		http.Error(w, dataErr.PublicMessage(), dataErr.HTTPStatus())
		return
	}
	http.Error(w, "file service failed", http.StatusInternalServerError)
}

func addDataPlaneCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeAssetHeaders(w http.ResponseWriter, contentType string, byteSize int64) {
	w.Header().Set("Content-Type", strings.TrimSpace(contentType))
	w.Header().Set("Content-Length", strconv.FormatInt(byteSize, 10))
}
