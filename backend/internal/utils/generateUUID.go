package utils

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
)

func GenerateUUID() (string, error) {
	val, err := uuid.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("error while generating UUID: %v", err)
	}
	return val, err
}
