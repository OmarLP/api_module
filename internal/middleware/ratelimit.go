package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/OmarLP/api_module/pkg/utils"
	"golang.org/x/time/rate"
)

type ipVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*ipVisitor)
	mu       sync.Mutex
)

// limpiar automaticamente los clientes inactivos cada 3 minutos para no saturar la memoria
func init() {
	go cleanVisitor()
}

func cleanVisitor() {
	for {
		time.Sleep(3 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// obtener el visitador
func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		// r = n° de eventos por segundo -- b = capacidad max en ráfaga (burst)
		limiter := rate.NewLimiter(r, b)
		visitors[ip] = &ipVisitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

// limitar las peticiones por ip
func RateLimitMiddleware(limit rate.Limit, burst int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// extraer ip cliente
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			limiter := getVisitor(ip, limit, burst)
			if !limiter.Allow() {
				utils.RespondError(w, http.StatusTooManyRequests, "too many requests, please try again later")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
