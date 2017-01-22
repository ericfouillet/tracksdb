package model

import (
	"crypto/rand"
	"fmt"
)

func getNextID(key string) (string, error) {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	res := int((int(b[0]) << 8) | int(b[1]))
	return fmt.Sprintf("%v-%x", key, res), nil
}
