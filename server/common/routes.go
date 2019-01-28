package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// SetupRoutes sets up all endpoints for the we service
func (s *System) SetupRoutes() {
	http.HandleFunc("/travel/v1/login", s.login)
	http.HandleFunc("/travel/v1/register", s.register)
}

func (s *System) login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type loginRequestBody struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	var rb loginRequestBody
	if err := json.Unmarshal(body, &rb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(rb.UserName) == 0 || len(rb.Password) == 0 {
		http.Error(w, "Missing request parameters", http.StatusBadRequest)
		return
	}

	if _, err := s.users.CheckPassword(rb.UserName, rb.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Create session cookie

	w.WriteHeader(200)

}

func (s *System) register(w http.ResponseWriter, r *http.Request) {

}
