package main

import (
	"syscall"
)

type netDevice struct {
	name       string
	macaddr    [6]uint8
	socket     int
	sockaddr   syscall.SockaddrLinklayer
	etheHeader ethernetHeader
}
