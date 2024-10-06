package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	byteRead, byteReadError := rand.Read(b)

	if byteReadError != nil {
		return nil, fmt.Errorf("byte reading error %w", byteReadError)
	}

	if byteRead < n {
		return nil, fmt.Errorf("Did not read enough byte")
	}

	return b, nil
}

func String(n int) (string, error) {
	b, byteError := Bytes(n)

	if byteError != nil {
		return "", byteError
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
