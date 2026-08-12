package util

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shohann/golang-ecommerce-api/apperr"
)

func SendData(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func SendError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(msg)
}

func SendAppError(w http.ResponseWriter, err error) {
	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		SendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch appErr.Kind {
	case apperr.KindValidation:
		SendError(w, http.StatusBadRequest, appErr.Message)
	case apperr.KindUnauthorized:
		SendError(w, http.StatusUnauthorized, appErr.Message)
	case apperr.KindNotFound:
		SendError(w, http.StatusNotFound, appErr.Message)
	case apperr.KindConflict:
		SendError(w, http.StatusConflict, appErr.Message)
	default:
		SendError(w, http.StatusInternalServerError, "internal server error")
	}
}
