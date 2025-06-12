package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"log/slog"
	"os"
	"github.com/rs/cors"

	"sync"
	"time"
	"golang.org/x/time/rate"
)

type CalcRequest struct {
	A  int    `json:"a"`
	B  int    `json:"b"`
	Op string `json:"op"`
}

type CalcRequest2 struct {
	A  int    `json:"a"`
	B  int    `json:"b"`
}

type CalcResponse struct {
	Result int `json:"result"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		logger.Error("Invalid method", "method", r.Method, "path", r.URL.Path)
		return
	}

	var req CalcRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logger.Error("Failed to decode JSON", "error", err, "path", r.URL.Path)
		return
	}

	// 👉 TODO: Implement operation logic here
	if req.Op == "" {
		http.Error(w, "Operation not specified", http.StatusBadRequest)
		logger.Error("operation not specified", "path", r.URL.Path)
		return
	}
	if req.Op != "add" && req.Op != "subtract" && req.Op != "multiply" && req.Op != "divide" {
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		logger.Error("invalid operator", "op", req.Op)
		return
	}
	var output int
	switch req.Op {
		case "add":
			output = req.A + req.B
		case "subtract":
			output = req.A - req.B
		case "multiply":
			output = req.A * req.B
		case "divide":
			if req.B == 0 {
				http.Error(w, "Division by zero", http.StatusBadRequest)
				logger.Error("division by zero", "a", req.A, "b", req.B)
				return
			}
			output = req.A / req.B
		default:
			http.Error(w, "Unknown operation", http.StatusBadRequest)
			logger.Error("invalid operator", "op", req.Op)
			return
	}

	res := CalcResponse{Result: output} // placeholder

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CalcRequest2
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	res := CalcResponse{Result: req.A + req.B}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CalcRequest2
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	res := CalcResponse{Result: req.A - req.B}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CalcRequest2
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	res := CalcResponse{Result: req.A * req.B}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CalcRequest2
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.B == 0 {
		http.Error(w, "Division by zero", http.StatusBadRequest)
		return
	}


	res := CalcResponse{Result: req.A / req.B}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}



var logger *slog.Logger



type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex // mutex to protect visitors map

func init() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
	log.Fatal(err)
}

logger = slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo}))

go cleanupVisitors()

}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()  // ensure thread-safe access to visitors map

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(0.2, 5) // 1 req per 5s, burst of 5
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			ip = fwd
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			logger.Warn("rate limit exceeded", "ip", ip)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}



func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		} else if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" {
			ip = fwdFor
		}

		logger.Info("received request",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", ip,
		)

		next.ServeHTTP(w, r)
	})
}


func main() {
	// http.HandleFunc("/calculate", calculateHandler)
	// fmt.Println("Server started at http://localhost:8080")
	// http.ListenAndServe(":8080", nil)
	server := http.NewServeMux()

	server.HandleFunc("/calculate", calculateHandler)
	server.HandleFunc("/add",addHandler)
	server.HandleFunc("/subtract", subtractHandler)
	server.HandleFunc("/multiply", multiplyHandler)
	server.HandleFunc("/divide", divideHandler)

	corshandler := cors.AllowAll().Handler(server)
	logHandler := loggingMiddleware(corshandler)
	rateLimitedHandler := rateLimitMiddleware(logHandler)
	
	fmt.Println("Server started at http://localhost:8080")
	logger.Info("Server starting on :8080")

	http.ListenAndServe(":8080", rateLimitedHandler)

	//test curl
	// curl -X GET http://localhost:8080/calculate -H "Content-Type: application/json" -d "{\"a\":10, \"b\":4, \"op\":\"add\"}"
	// curl -X GET http://localhost:8080/add -H "Content-Type: application/json" -d "{\"a\":10, \"b\":4}"



}
