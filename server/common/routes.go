package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// SetupRoutes sets up all endpoints for the we service
func (s *System) SetupRoutes() {
	http.HandleFunc("/", s.base)
	http.HandleFunc("/set-value", s.setValue)
	http.HandleFunc("/get-value", s.getValue)
}

func (s *System) setValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST request implemented", http.StatusNotImplemented)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	type requestBodyType struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	var rb requestBodyType
	if err := json.Unmarshal(body, &rb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.kvStore.Store(rb.Key, rb.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

func (s *System) getValue(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	key := keys.Get("key")
	if key == "" {
		http.Error(w, "No key query parameter found", http.StatusBadRequest)
		return
	}
	value, err := s.kvStore.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(*value))
}

func (s *System) base(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Welcome to a Go web server with key value storage"))
}
