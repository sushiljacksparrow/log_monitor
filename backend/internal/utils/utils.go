package utils

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/hashicorp/go-uuid"
)

// GenerateUUID creates a unique identifier
func GenerateUUID() (string, error) {
	val, err := uuid.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("error while generating UUID: %v", err)
	}
	return val, err
}

// RandomIP generates a random IPv4 address
func RandomIP() net.IP {
	rand.Seed(time.Now().UnixNano())

	return net.IPv4(
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	)
}
