package main

import (
	"os"
	"strings"
)

func GetHostname() string {
	var host_re = "Unknown"
	hostname, err := os.Hostname()
	if err != nil {
		return host_re
	}
	host_re = strings.ToLower(hostname)
	return host_re
}