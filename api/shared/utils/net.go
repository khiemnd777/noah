package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// CheckPortOpen returns true if host:port is reachable via TCP
func CheckPortOpen(host string, port int) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// DetectPIDFromPort attempts to find PID listening on given port (cross-platform)
func DetectPIDFromPort(port int) (int, error) {
	switch runtime.GOOS {
	case "windows":
		return detectPIDWindows(port)
	case "linux":
		return detectPIDLinux(port)
	case "darwin":
		return detectPIDDarwin(port)
	default:
		return 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func detectPIDDarwin(port int) (int, error) {
	cmd := exec.Command("lsof", "-nP", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN", "-t")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("lsof failed: %w", err)
	}

	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return 0, errors.New("no PID found")
	}
	return strconv.Atoi(fields[0])
}

func detectPIDLinux(port int) (int, error) {
	inodes := map[string]struct{}{}
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		found, err := listeningSocketInodes(path, port)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return 0, err
		}
		for inode := range found {
			inodes[inode] = struct{}{}
		}
	}
	if len(inodes) == 0 {
		return 0, fmt.Errorf("no listening socket found for port %d", port)
	}

	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return 0, fmt.Errorf("read /proc failed: %w", err)
	}
	for _, procEntry := range procEntries {
		if !procEntry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(procEntry.Name())
		if err != nil {
			continue
		}

		fdDir := filepath.Join("/proc", procEntry.Name(), "fd")
		fdEntries, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fdEntry := range fdEntries {
			target, err := os.Readlink(filepath.Join(fdDir, fdEntry.Name()))
			if err != nil {
				continue
			}
			if !strings.HasPrefix(target, "socket:[") || !strings.HasSuffix(target, "]") {
				continue
			}
			inode := strings.TrimSuffix(strings.TrimPrefix(target, "socket:["), "]")
			if _, ok := inodes[inode]; ok {
				return pid, nil
			}
		}
	}

	return 0, fmt.Errorf("no PID found for listening port %d", port)
}

func listeningSocketInodes(path string, port int) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	inodes := map[string]struct{}{}
	targetPort := strings.ToUpper(fmt.Sprintf("%04X", port))
	lines := strings.Split(string(data), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		local := fields[1]
		state := fields[3]
		inode := fields[9]
		if state != "0A" {
			continue
		}
		parts := strings.Split(local, ":")
		if len(parts) != 2 {
			continue
		}
		if strings.EqualFold(parts[1], targetPort) {
			inodes[inode] = struct{}{}
		}
	}

	return inodes, nil
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
				pidStr := fields[len(fields)-1]
				return strconv.Atoi(pidStr)
			}
		}
	}
	return 0, fmt.Errorf("no matching PID found for port %d", port)
}
