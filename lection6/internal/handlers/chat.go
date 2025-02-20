package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode/utf8"

	"entrance/lection6/internal/middlewares"
	"entrance/lection6/internal/models"
	"entrance/lection6/internal/storage"

	"github.com/go-chi/chi/v5"
)

const (
	defaultMaxPublicMessageLength  = 1023
	defaultMaxPrivateMessageLength = 4095
)

type ChatService struct {
	repo storage.Repository
	Configs
}

type Configs struct {
	maxPublicMessageLength  int
	maxPrivateMessageLength int
}

func DefaultChatService(repo storage.Repository) *ChatService {
	return &ChatService{
		repo: repo,
		Configs: Configs{
			maxPublicMessageLength:  defaultMaxPublicMessageLength,
			maxPrivateMessageLength: defaultMaxPrivateMessageLength,
		},
	}
}

func NewChatService(repo storage.Repository, configs Configs) *ChatService {
	return &ChatService{repo: repo, Configs: configs}
}

func (s *ChatService) ReadPublic(w http.ResponseWriter, r *http.Request) {
	var chatName = chi.URLParam(r, "chatName")
	messages := s.repo.GetPublicMessages(chatName)
	if len(messages) == 0 {
		http.Error(w, fmt.Sprintf("Public chat %s not found", chatName), http.StatusNotFound)
	}

	bytes, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func (s *ChatService) ReadPrivate(w http.ResponseWriter, r *http.Request) {
	chatName := chi.URLParam(r, "chatName")
	userName := r.Context().Value(middlewares.UserName).(string)
	messages := s.repo.GetPrivateMessages(userName, chatName)
	bytes, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func (s *ChatService) SendToPublic(w http.ResponseWriter, r *http.Request) {
	s.sendMessage(w, r, s.maxPublicMessageLength, s.repo.AddPublicMessage)
}

func (s *ChatService) ListPublic(w http.ResponseWriter, _ *http.Request) {
	chats := s.repo.GetAllPublicChats()
	bytes, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func (s *ChatService) SendPrivate(w http.ResponseWriter, r *http.Request) {
	s.sendMessage(w, r, s.maxPrivateMessageLength, s.repo.AddPrivateMessage)
}

func (s *ChatService) ListPrivate(w http.ResponseWriter, r *http.Request) {
	userName := r.Context().Value(middlewares.UserName).(string)
	chats := s.repo.GetAllPrivateChats(userName)
	bytes, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func (s *ChatService) sendMessage(
	w http.ResponseWriter, r *http.Request, maxMessageLength int, addMessageFunc func(string, models.Message)) {
	limitedReader := io.LimitReader(r.Body, int64(maxMessageLength)+1)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if utf8.RuneCount(body) > maxMessageLength {
		http.Error(
			w, fmt.Sprintf("Message larger than %d characters", maxMessageLength),
			http.StatusRequestEntityTooLarge,
		)
		return
	}

	msg := models.Message{
		Sender:  r.Context().Value(middlewares.UserName).(string),
		Time:    time.Now(),
		Content: string(body),
	}

	addMessageFunc(chi.URLParam(r, "chatName"), msg)
}
