package network

import (
	"math/rand"
	"net"
	"time"
)

func RandomIP() net.IP {
	rand.Seed(time.Now().UnixNano())

	return net.IPv4(
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	)
}
