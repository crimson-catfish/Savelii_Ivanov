package reopositories

import (
	"entrance/lection6/internal/models"
)

type Repository interface {
	UserExists(name string) bool
	AddUser(credentials models.Credentials)
	GetPassword(name string) string

	GetAllPublicChats() []string
	GetPublicMessages(chat string) []models.Message
	AddPublicMessage(chat string, msg models.Message)

	GetAllPrivateChats(user string) []string
	GetPrivateMessages(userName, chatName string) []models.Message
	AddPrivateMessage(receiver string, msg models.Message)
}
