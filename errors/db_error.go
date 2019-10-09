package errors

import (
	"encoding/json"
	"strings"
)

func IsRecordNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "record not found")
}

func IsDuplicateKeyError(err error) (string, bool) {
	ok := strings.Contains(err.Error(), "duplicate key")

	b, err := json.Marshal(err)
	if err != nil {
		return "", ok
	}

	errMsg := struct {
		Detail string `json:"detail"`
	}{}

	if err := json.Unmarshal(b, &errMsg); err != nil {
		return "", ok
	}

	return errMsg.Detail, ok
}
