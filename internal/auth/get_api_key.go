package auth

import (
	"errors"
	"net/http"
	"strings"
)

func (a *Auth) GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")

	if apiKey == "" {
		return "", errors.New("no authorizarion headers found")

	}

	key := strings.Fields(apiKey)
	if len(key) < 2 {
		return "", errors.New("no api key found in headers")
	}

	if key[0] != "ApiKey" {
		return "", errors.New("no api key found in headers")
	}

	return key[1], nil

}
