package runtime

import (
	"fmt"
	"net"
)

func reservedListener(host string, port int) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %w", addr, err)
	}
	return listener, nil
}
