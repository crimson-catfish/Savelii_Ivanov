package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode/utf8"

	"entrance/lection5/database"
	"entrance/lection5/middlewares"
	"entrance/lection5/models"

	"github.com/go-chi/chi/v5"
)

const (
	maxPublicMessageLength  = 1023
	maxPrivateMessageLength = 4095
)

func ReadPublicChat(w http.ResponseWriter, r *http.Request) {
	var chatName = chi.URLParam(r, "chatName")
	messages := database.GetPublicMessages(chatName)
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

func ReadPrivateChat(w http.ResponseWriter, r *http.Request) {
	chatName := chi.URLParam(r, "chatName")
	userName := r.Context().Value(middlewares.UserName).(string)
	messages := database.GetPrivateMessages(userName, chatName)
	bytes, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func SendToPublicChat(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, r, maxPublicMessageLength, database.AddPublicMessage)
}

func ListPublicChats(w http.ResponseWriter, _ *http.Request) {
	chats := database.GetAllPublicChats()
	bytes, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func SendPrivate(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, r, maxPrivateMessageLength, database.AddPrivateMessage)
}

func ListPrivateChats(w http.ResponseWriter, r *http.Request) {
	userName := r.Context().Value(middlewares.UserName).(string)
	chats := database.GetAllPrivateChats(userName)
	bytes, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func sendMessage(
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
