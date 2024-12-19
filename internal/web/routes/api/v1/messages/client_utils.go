package messages

import (
	"encoding/json"

	"github.com/k10wl/hermes/internal/validator"
)

func typeDetector(data []byte) (string, error) {
	type typeDetector struct {
		Type string `json:"type,required"`
	}
	var t typeDetector
	err := json.Unmarshal(data, &t)
	return t.Type, err
}

func decode(receiver any, data []byte) error {
	err := json.Unmarshal(data, receiver)
	if err != nil {
		return err
	}
	err = validator.Validate.Struct(receiver)
	return err
}
