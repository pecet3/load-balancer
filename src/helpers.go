package main

import (
	"net"
	"net/http"
	"strings"
)

func getIP(r *http.Request) string {
	xri := r.Header.Get("X-Real-Ip")
	if xri != "" {
		ip := xri
		if idx := net.ParseIP(ip); idx != nil {
			return ip
		}
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip := xff
		if idx := net.ParseIP(ip); idx != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ""
	}
	return userIP.String()
}
func isWebSocketRequest(r *http.Request) bool {
	upgrade := strings.ToLower(r.Header.Get("Upgrade"))
	connection := strings.ToLower(r.Header.Get("Connection"))

	return upgrade == "websocket" && strings.Contains(connection, "upgrade")
}
