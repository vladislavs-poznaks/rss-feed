package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts API key from headers in HTTP request
// Example:
// Authorization: ApiKey {api_key_here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authorization header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("incorrect authorization header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("incorrect authorization header")
	}

	return vals[1], nil
}
