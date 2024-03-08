package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var cache *LRUCache = NewLRUCache(1024)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if value, ok := cache.Get(key); ok {
		json.NewEncoder(w).Encode(map[string]interface{}{"value": value})
	} else {
		http.Error(w, "Key not found", http.StatusNotFound)
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	var value interface{}
	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Assuming expiration is passed in seconds for simplicity
	durationStr := r.URL.Query().Get("expiration")
	duration, _ := time.ParseDuration(durationStr + "s")
	cache.Set(key, value, duration)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/get", enableCORS(getHandler))
	http.HandleFunc("/set", enableCORS(setHandler))
	http.ListenAndServe(":8080", nil)
}
