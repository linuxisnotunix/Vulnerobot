package tools

import (
	"net"
	"time"
)

//IsTCPAccessible try to open a TCP socket and return true if success
func IsTCPAccessible(host string) (bool, error) {
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}
