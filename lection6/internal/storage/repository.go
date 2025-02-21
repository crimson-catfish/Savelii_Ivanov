package storage

import (
	"entrance/lection6/internal/models"
)

type Repository interface {
	UserExists(name string) (bool, error)
	AddUser(credentials models.Credentials) error
	GetPassword(name string) (string, error)

	GetAllPublicChats() ([]string, error)
	GetPublicMessages(chat string) ([]models.Message, error)
	AddPublicMessage(chat string, msg models.Message) error

	GetAllPrivateChats(user string) ([]string, error)
	GetPrivateMessages(userName, chatName string) ([]models.Message, error)
	AddPrivateMessage(receiver string, msg models.Message) error

	Close() error
}
