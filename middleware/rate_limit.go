package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// User Rate Limiter
type UserRateLimiter struct {
	users map[string]*rate.Limiter
	mu    *sync.RWMutex
	r     rate.Limit
	b     int
}

func NewUserRateLimiter(r rate.Limit, b int) *UserRateLimiter {
	return &UserRateLimiter{
		users: make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		r:     r,
		b:     b,
	}
}

func (u *UserRateLimiter) GetLimiter(userID string) *rate.Limiter {
	u.mu.Lock()
	defer u.mu.Unlock()

	limiter, exists := u.users[userID]
	if !exists {
		limiter = rate.NewLimiter(u.r, u.b)
		u.users[userID] = limiter
	}

	return limiter
}

var loginLimiter = NewIPRateLimiter(rate.Every(1*time.Second/3), 3)

var bookingLimiter = NewUserRateLimiter(rate.Every(1*time.Second/5), 5)

func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func LoginRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := loginLimiter.GetLimiter(ip)
		
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		
		next(w, r)
	}
}

func BookingRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(string)
		if !ok {
			ip := getIP(r)
			limiter := loginLimiter.GetLimiter(ip)
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		} else {
			limiter := bookingLimiter.GetLimiter(userID)
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}
		
		next(w, r)
	}
}
