package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os/exec"
	goruntime "runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Reserved struct {
	Listener net.Listener
	Port     int
}

func FindAvailablePort(host string) (*Reserved, error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if err != nil {
		if listener != nil {
			_ = listener.Close()
		}
		if opErr, ok := err.(*net.OpError); ok && opErr.Err == syscall.EADDRNOTAVAIL {
			return nil, fmt.Errorf("host %s not reachable: %w", host, err)
		}
		return nil, fmt.Errorf("cannot allocate port: %w", err)
	}
	return &Reserved{
		Listener: listener,
		Port:     listener.Addr().(*net.TCPAddr).Port,
	}, nil
}

func ListenOnAvailablePort(host string, port int) (*Reserved, error) {
	if port == 0 {
		return FindAvailablePort(host)
	}

	tryBind := func(p int) (net.Listener, error) {
		return net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(p)))
	}

	listener, err := tryBind(port)
	if err == nil {
		return &Reserved{Listener: listener, Port: port}, nil
	}

	if errors.Is(err, syscall.EADDRINUSE) || errors.Is(err, syscall.EACCES) {
		return FindAvailablePort(host)
	}

	return nil, fmt.Errorf("listen failed: %w", err)
}

func EnsurePortAvailable(host string, desiredPort int) (*Reserved, error) {
	if desiredPort == 0 {
		return FindAvailablePort(host)
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(desiredPort)))
	if err == nil {
		return &Reserved{Listener: listener, Port: desiredPort}, nil
	}
	if listener != nil {
		_ = listener.Close()
	}

	if errors.Is(err, syscall.EADDRINUSE) || errors.Is(err, syscall.EACCES) {
		return FindAvailablePort(host)
	}

	return nil, fmt.Errorf("cannot bind desired port %d: %w", desiredPort, err)
}

func CheckPortOpen(host string, port int) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func DetectPIDFromPort(port int) (int, error) {
	switch goruntime.GOOS {
	case "windows":
		return detectPIDWindows(port)
	case "linux", "darwin":
		return detectPIDUnix(port)
	default:
		return 0, fmt.Errorf("unsupported OS: %s", goruntime.GOOS)
	}
}

func detectPIDUnix(port int) (int, error) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("lsof failed: %w", err)
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return 0, errors.New("no PID found")
	}
	return strconv.Atoi(pidStr)
}

func detectPIDWindows(port int) (int, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("netstat failed: %w", err)
	}

	lines := bytes.Split(output, []byte("\n"))
	target := fmt.Sprintf(":%d", port)
	for _, line := range lines {
		str := string(line)
		if strings.Contains(str, target) && strings.Contains(str, "LISTENING") {
			fields := strings.Fields(str)
			if len(fields) >= 5 {
				return strconv.Atoi(fields[len(fields)-1])
			}
		}
	}
	return 0, fmt.Errorf("no matching PID found for port %d", port)
}
