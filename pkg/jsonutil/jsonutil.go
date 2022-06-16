package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeJSON(r io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}
	r.Close()

	return nil
}

func EncodeJSON(v interface{}, w io.Writer) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("unable to json encode: %w", err)
	}

	return nil
}
