package rand

import (
	"crypto/rand"
	"fmt"
)

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
