package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"io"
	"log"
	"net/http"
	"strings"
)

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Could not write error message [%s]", err.Error())
		panic(err)
	}
}

func writeJson(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshalling json response [%v]", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type malformedRequestError struct {
	Status  int
	Message string
}

func (mr *malformedRequestError) Error() string {
	return mr.Message
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			Message := "Content-Type header is not application/json"
			return &malformedRequestError{Status: http.StatusUnsupportedMediaType, Message: Message}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			Message := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}

		case errors.Is(err, io.ErrUnexpectedEOF):
			Message := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}

		case errors.As(err, &unmarshalTypeError):
			Message := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			Message := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}

		case errors.Is(err, io.EOF):
			Message := "Request body must not be empty"
			return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}

		case err.Error() == "http: request body too large":
			Message := "Request body must not be larger than 1MB"
			return &malformedRequestError{Status: http.StatusRequestEntityTooLarge, Message: Message}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		Message := "Request body must only contain a single JSON object"
		return &malformedRequestError{Status: http.StatusBadRequest, Message: Message}
	}

	return nil
}
