package utils

import frameworkruntime "github.com/khiemnd777/noah_framework/runtime"

func CheckPortOpen(host string, port int) bool {
	return frameworkruntime.CheckPortOpen(host, port)
}

func DetectPIDFromPort(port int) (int, error) {
	return frameworkruntime.DetectPIDFromPort(port)
}
