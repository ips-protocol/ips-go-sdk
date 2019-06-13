package netools

import (
	"net"
	"strconv"
)

func GetFreePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()

	port = l.Addr().(*net.TCPAddr).Port
	return
}

func IsLocalPortAvailable(port int) (available bool, err error) {
	lport := "127.0.0.1:" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", lport)
	if err != nil {
		return
	}
	available = true
	err = ln.Close()
	return
}
