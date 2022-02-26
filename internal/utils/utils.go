package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

//go:generate mockgen -destination=validator_mock.go -package=utils . Validator
type Validator interface {
	IsValid() bool
}

func ValidateBody(b []byte, v Validator) (int, error) {
	err := json.Unmarshal(b, &v)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !v.IsValid() {
		return http.StatusBadRequest, errors.New("invalid body")
	}
	return 0, nil
}

// LoadEnv wraps os.Getenv and returns a default value if env is empty
func LoadEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

// String returns reference to a string value
func String(s string) *string {
	return &s
}
