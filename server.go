package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	
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
		return
	}

	var req CalcRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 👉 TODO: Implement operation logic here
	if req.Op == "" {
		http.Error(w, "Operation not specified", http.StatusBadRequest)
		return
	}
	if req.Op != "add" && req.Op != "subtract" && req.Op != "multiply" && req.Op != "divide" {
		http.Error(w, "Invalid operation", http.StatusBadRequest)
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
				return
			}
			output = req.A / req.B
		default:
			http.Error(w, "Unknown operation", http.StatusBadRequest)
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
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", server)

	//test curl
	// curl -X GET http://localhost:8080/calculate -H "Content-Type: application/json" -d "{\"a\":10, \"b\":4, \"op\":\"add\"}"
	// curl -X GET http://localhost:8080/add -H "Content-Type: application/json" -d "{\"a\":10, \"b\":4}"



}
