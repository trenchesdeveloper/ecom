package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ParseJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// limit the request body to 1MB
	maxBytes := 1024 * 1024

	// read the request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// decode the request body into the data interface
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("request body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			}
			return fmt.Errorf("request body contains an invalid value (at position %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown field %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("request body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling request body: %s", err.Error())
		default:
			return err
		}
	}

	// check if it's not more than one json object
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("the request body must only contain a single JSON object")
	}

	return nil

}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	//
	if headers != nil {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	// set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// write the status code
	w.WriteHeader(statusCode)

	// write the JSON response
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON writes an error to the response
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusInternalServerError

	if status != nil {
		statusCode = status[0]
	}

	var payload JSONResponse

	payload.Error = true
	payload.Message = err.Error()

	WriteJSON(w, statusCode, payload)

}
