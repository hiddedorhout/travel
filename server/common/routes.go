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
	http.HandleFunc("/travel/v1/travel-session", s.validateAuthorization(s.travelHandler))
	http.HandleFunc("/travel/v1/routes-list", s.validateAuthorization(s.getRoutesList))
	http.HandleFunc("/travel/v1/route", s.validateAuthorization(s.getRoute))
}

func (s *System) login(w http.ResponseWriter, r *http.Request) {
	type loginRequestBody struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rb loginRequestBody
	if err := json.Unmarshal(body, &rb); err != nil {
		errorResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(rb.UserName) == 0 || len(rb.Password) == 0 {
		errorResponseWriter(w, "Missing request parameters", http.StatusBadRequest)
		return
	}

	id, err := s.users.CheckPassword(rb.UserName, rb.Password)
	if err != nil {
		errorResponseWriter(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Create session cookie
	session, err := s.sessions.GenerateSession(*id)
	if err != nil {
		errorResponseWriter(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(*session))

}

func (s *System) register(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	type responseBody struct {
		UserID string `json:"userID"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rb requestBody
	if err := json.Unmarshal(body, &rb); err != nil {
		errorResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(rb.UserName) == 0 || len(rb.Password) == 0 {
		errorResponseWriter(w, "Missing request body parameters", http.StatusBadRequest)
		return
	}
	id, err := s.users.RegisterUser(rb.UserName, rb.Password)
	if err != nil {
		errorResponseWriter(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respb := responseBody{
		UserID: *id,
	}

	resp, err := json.Marshal(respb)
	if err != nil {
		errorResponseWriter(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *System) travelHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		s.startTravelSession(w, r)
	case "PUT":
		s.updateTravelSession(w, r)
	case "GET":
		s.getTravelSession(w, r)
	case "DELETE":
		s.deleteTravelSession(w, r)
	default:
		errorResponseWriter(w, "Unsuported method", http.StatusNotImplemented)
		return
	}
}

func (s *System) startTravelSession(w http.ResponseWriter, r *http.Request) {

}

func (s *System) updateTravelSession(w http.ResponseWriter, r *http.Request) {

}

func (s *System) getTravelSession(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (s *System) deleteTravelSession(w http.ResponseWriter, r *http.Request) {

}

func (s *System) getRoutesList(w http.ResponseWriter, r *http.Request) {

}

func (s *System) routesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getRoute(w, r)
	case "POST":
		s.createRoute(w, r)
	default:
		errorResponseWriter(w, "Unsuported method", http.StatusNotImplemented)
		return
	}
}

func (s *System) getRoute(w http.ResponseWriter, r *http.Request) {

}

func (s *System) createRoute(w http.ResponseWriter, r *http.Request) {

}

// UTILS
func errorResponseWriter(w http.ResponseWriter, message string, responseCode int) {
	type errorResponse struct {
		Reason string `json:"reason"`
	}
	rb := errorResponse{
		Reason: message,
	}
	resp, _ := json.Marshal(rb)
	w.WriteHeader(responseCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *System) validateAuthorization(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			errorResponseWriter(w, "Missing Authorization header", http.StatusBadRequest)
			return
		}
		if err := s.sessions.ValidateSession(auth); err != nil {
			errorResponseWriter(w, err.Error(), http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}
