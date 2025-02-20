package database

import (
	"time"

	"entrance/lection6/internal/models"
)

// mock database populated with some sample data
var (
	userPasswords = map[string]string{"bill": "bill_rules", "uma_thruman": "bill_sucks"}

	publicChatsSlice = []string{"greatest-chat-of-all-times", "ok-chat", "meh-chat"}

	publicChatsMap = map[string][]models.Message{
		"greatest-chat-of-all-times": make([]models.Message, 0),
		"ok-chat":                    make([]models.Message, 0),
		"meh-chat":                   make([]models.Message, 0),
	}

	privateChatsSlice = []models.PrivateChat{
		{
			User1: "bill",
			User2: "uma_thruman",
			Messages: []models.Message{
				{
					Sender:  "uma_thruman",
					Time:    time.Date(2003, 10, 10, 0, 0, 0, 0, time.UTC),
					Content: "ima kill ya",
				},
			},
		},
	}

	privateChatsMap = map[string]map[string]*models.PrivateChat{
		"bill":        {"uma_thruman": &privateChatsSlice[0]},
		"uma_thruman": {"bill": &privateChatsSlice[0]},
	}
)

type MockRepository struct{}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (repo *MockRepository) UserExists(name string) bool {
	_, ok := userPasswords[name]
	return ok
}

func (repo *MockRepository) AddUser(credentials models.Credentials) {
	userPasswords[credentials.Name] = credentials.Password
}

func (repo *MockRepository) GetPassword(name string) string {
	return userPasswords[name]
}

func (repo *MockRepository) AddPublicMessage(chat string, msg models.Message) {
	if publicChatsMap[chat] == nil {
		publicChatsMap[chat] = make([]models.Message, 0)
		publicChatsSlice = append(publicChatsSlice, chat)
	}

	publicChatsMap[chat] = append(publicChatsMap[chat], msg)
}

func (repo *MockRepository) GetPublicMessages(chat string) []models.Message {
	return publicChatsMap[chat]
}

func (repo *MockRepository) GetAllPublicChats() []string {
	return publicChatsSlice
}

func (repo *MockRepository) AddPrivateMessage(receiver string, msg models.Message) {
	if privateChatsMap[receiver] == nil {
		privateChatsMap[receiver] = make(map[string]*models.PrivateChat)
	}
	if privateChatsMap[msg.Sender] == nil {
		privateChatsMap[msg.Sender] = make(map[string]*models.PrivateChat)
	}

	if privateChatsMap[receiver][msg.Sender] == nil {
		chat := models.PrivateChat{
			User1:    msg.Sender,
			User2:    receiver,
			Messages: make([]models.Message, 0),
		}
		privateChatsSlice = append(privateChatsSlice, chat)

		chatPtr := &privateChatsSlice[len(privateChatsSlice)-1]
		privateChatsMap[receiver][msg.Sender] = chatPtr
		privateChatsMap[msg.Sender][receiver] = chatPtr
	}

	privateChatsMap[receiver][msg.Sender].Messages = append(privateChatsMap[receiver][msg.Sender].Messages, msg)
}

func (repo *MockRepository) GetAllPrivateChats(user string) []string {
	chats := make([]string, 0)
	for chat := range privateChatsMap[user] {
		chats = append(chats, chat)
	}
	return chats
}

func (repo *MockRepository) GetPrivateMessages(userName, chatName string) []models.Message {
	if privateChatsMap[userName] == nil || privateChatsMap[userName][chatName] == nil {
		return []models.Message{}
	}
	return privateChatsMap[userName][chatName].Messages
}
