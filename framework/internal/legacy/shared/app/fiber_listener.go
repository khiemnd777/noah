package app

import frameworkruntime "github.com/khiemnd777/noah_framework/runtime"

type Reserved = frameworkruntime.Reserved

func ListenOnAvailablePort(host string, port int) (*Reserved, error) {
	return frameworkruntime.ListenOnAvailablePort(host, port)
}

func FindAvailablePort(host string) (*Reserved, error) {
	return frameworkruntime.FindAvailablePort(host)
}

func EnsurePortAvailable(host string, desiredPort int) (*Reserved, error) {
	return frameworkruntime.EnsurePortAvailable(host, desiredPort)
}
