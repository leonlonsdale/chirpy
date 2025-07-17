package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrorNoBearerToken = errors.New("no bearer token in headers")

func GetBearerToken(headers http.Header) (string, error) {

	authorisation := headers.Get("Authorization")
	if authorisation == "" {
		return "", ErrorNoBearerToken
	}

	token := strings.Fields(authorisation)

	if len(token) < 2 {
		return "", ErrorNoBearerToken
	}

	if token[0] != "Bearer" {
		return "", ErrorNoBearerToken
	}

	return token[1], nil
}
