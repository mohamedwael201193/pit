package httpx

import (
	"net"
	"net/http"
	"strings"
)

func HealthOriginOK(origin string) bool {
	o := strings.TrimSpace(origin)
	if o == "" {
		return false
	}
	if o == "https://pit0g.vercel.app" {
		return true
	}
	if strings.HasPrefix(o, "https://") && strings.HasSuffix(o, ".vercel.app") {
		return true
	}
	if loopbackHTTP(o) {
		return true
	}
	return o == "https://tauri.localhost" || o == "tauri://localhost"
}

func CompanionOriginOK(origin string) bool {
	switch strings.TrimSpace(origin) {
	case "https://pit0g.vercel.app",
		"http://127.0.0.1:4173", "http://localhost:4173",
		"http://127.0.0.1:5173", "http://localhost:5173",
		"http://127.0.0.1:3000", "http://localhost:3000",
		"http://127.0.0.1:3001", "http://localhost:3001",
		"https://tauri.localhost", "tauri://localhost":
		return true
	default:
		return false
	}
}

func CodeOriginOK(origin string) bool {
	switch strings.TrimSpace(origin) {
	case "",
		"http://127.0.0.1:3001", "http://localhost:3001",
		"https://tauri.localhost", "tauri://localhost":
		return true
	default:
		return false
	}
}

func loopbackHTTP(o string) bool {
	return strings.HasPrefix(o, "http://127.0.0.1:") || strings.HasPrefix(o, "http://localhost:") ||
		o == "http://127.0.0.1" || o == "http://localhost"
}

func IsLoopbackAddr(remote string) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		host = remote
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func Public(next http.Handler) http.Handler {
	return cors(HealthOriginOK, false, next)
}

func Companion(next http.Handler) http.Handler {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" && !CompanionOriginOK(origin) {
			http.Error(w, "origin_denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
	return loopbackOnly(cors(CompanionOriginOK, true, inner))
}

func loopbackOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr != "" && !IsLoopbackAddr(r.RemoteAddr) {
			http.Error(w, "loopback_only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func cors(ok func(string) bool, denyBadOrigin bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && ok(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}
		if r.Method == http.MethodOptions {
			if origin != "" && !ok(origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if denyBadOrigin && origin != "" && !ok(origin) {
			http.Error(w, "origin_denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
