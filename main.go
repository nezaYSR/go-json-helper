package json_helper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ReadJSON(r io.Reader, data interface{}) error {
	maxBytes := 1048576
	r = io.LimitReader(r, int64(maxBytes))

	d := json.NewDecoder(r)
	err := d.Decode(data)
	if err != nil {
		return err
	}

	err = d.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func WriteJSON(w io.Writer, status int, data interface{}, headers ...http.Header) error {
	o, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(o)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w io.Writer, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}
