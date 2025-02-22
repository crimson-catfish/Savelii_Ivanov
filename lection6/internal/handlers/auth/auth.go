package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"entrance/lection6/internal/models"
	"entrance/lection6/internal/storage"
	"entrance/lection6/pkg/auth"
)

type Service struct {
	repo storage.Repository
}

func NewAuthService(repo storage.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignIn(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println(err)
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

	exists, err := s.repo.UserExists(credentials.Name)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, fmt.Sprintf("User \"%s\" not found", credentials.Name), http.StatusConflict)
		return
	}

	password, err := s.repo.GetPassword(credentials.Name)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if password != credentials.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateJWTToken(credentials.Name)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (s *Service) SignUp(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println(err)
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

	exists, err := s.repo.UserExists(credentials.Name)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, fmt.Sprintf("Name \"%s\" is already taken", credentials.Name), http.StatusConflict)
		return
	}

	if err := s.repo.AddUser(credentials); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	token, err := auth.CreateJWTToken(credentials.Name)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}
