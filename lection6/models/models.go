package models

import "time"

type (
	Credentials struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	Message struct {
		Sender  string    `json:"sender"`
		Time    time.Time `json:"time"`
		Content string    `json:"content"`
	}

	PrivateChat struct {
		User1    string    `json:"user1"`
		User2    string    `json:"user2"`
		Messages []Message `json:"messages"`
	}
)
