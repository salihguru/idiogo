package xip

import (
	"net"
	"strings"
)

const (
	HeaderCloudflareIP = "CF-Connecting-IP"
	HeaderForwardedIP  = "X-Forwarded-For"
	HeaderRealIP       = "X-Real-IP"
)

type HeaderGetter func(name string) string

func ClaimRealIP(g HeaderGetter, def string) string {
	return claimIPs(g, def, claimCloudflareIP, claimForwardedIP, claimRealIP)
}

func claimCloudflareIP(g HeaderGetter) string {
	return g(HeaderCloudflareIP)
}

func claimForwardedIP(g HeaderGetter) string {
	return g(HeaderForwardedIP)
}

func claimRealIP(g HeaderGetter) string {
	return g(HeaderRealIP)
}

func claimIPs(g HeaderGetter, def string, claimers ...func(HeaderGetter) string) string {
	for _, claimer := range claimers {
		ip := claimer(g)
		if ip != "" {
			if strings.Contains(ip, ",") {
				ip = strings.Split(ip, ",")[0]
			}
			if isValidIP(ip) {
				return ip
			}
		}
	}
	return def
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return !parsedIP.IsPrivate() && !parsedIP.IsLoopback() && !parsedIP.IsUnspecified()
}
