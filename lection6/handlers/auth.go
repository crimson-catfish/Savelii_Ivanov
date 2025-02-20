package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"entrance/lection6/auth"
	"entrance/lection6/models"
	"entrance/lection6/reopositories"
)

type AuthService struct {
	repo reopositories.Repository
}

func NewAuthService(repo reopositories.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if credentials.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if credentials.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	if !s.repo.UserExists(credentials.Name) {
		http.Error(w, fmt.Sprintf("User \"%s\" not found", credentials.Name), http.StatusConflict)
		return
	}

	password := s.repo.GetPassword(credentials.Name)
	if password != credentials.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateJWTToken(credentials.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (s *AuthService) SignUp(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if credentials.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if credentials.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	if s.repo.UserExists(credentials.Name) {
		http.Error(w, fmt.Sprintf("Name \"%s\" is already taken", credentials.Name), http.StatusConflict)
		return
	}
	s.repo.AddUser(credentials)

	token, err := auth.CreateJWTToken(credentials.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}
