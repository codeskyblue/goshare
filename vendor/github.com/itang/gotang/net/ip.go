package net

import (
	"errors"
	"net"
	"os"
	"runtime"
	"strings"
)

// @TODO Improve
func LookupWlanIP4addr() (ip4 string, err error) {
	switch runtime.GOOS {
	case "linux":
		ifi, err := firstIfiNoErr([]string{"wlan0", "eth0"})
		addrs, err := ifi.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addr.String()
			if wlanIP4addrLike(ip) {
				end := len(ip) - 3
				return ip[0:end], nil
			}
		}
	default:
		name, err := os.Hostname()
		if err != nil {
			return "", err
		}
		addrs, err := net.LookupHost(name)
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if wlanIP4addrLike(addr) {
				return addr, nil
			}
		}
	}
	return "", errors.New("NO FOUND")
}

func firstIfiNoErr(names []string) (ifi *net.Interface, err error) {
	for _, name := range names {
		ifi, err = net.InterfaceByName(name)
		if err == nil {
			return ifi, nil
		}
	}
	return nil, err
}

func wlanIP4addrLike(ip string) bool {
	return !strings.Contains(ip, ":") && !strings.HasPrefix(ip, "127")
}
