package mock

import (
	"errors"
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

type Repository struct{}

func NewMockRepository() *Repository {
	return &Repository{}
}

func (repo *Repository) UserExists(name string) (bool, error) {
	_, ok := userPasswords[name]
	return ok, nil
}

func (repo *Repository) AddUser(credentials models.Credentials) error {
	_, exists := userPasswords[credentials.Name]
	if exists {
		return errors.New("user already exists")
	}
	userPasswords[credentials.Name] = credentials.Password
	return nil
}

func (repo *Repository) GetPassword(name string) (string, error) {
	password, exists := userPasswords[name]
	if !exists {
		return "", errors.New("user not found")
	}
	return password, nil
}

func (repo *Repository) AddPublicMessage(chat string, msg models.Message) error {
	if publicChatsMap[chat] == nil {
		publicChatsMap[chat] = make([]models.Message, 0)
		publicChatsSlice = append(publicChatsSlice, chat)
	}

	publicChatsMap[chat] = append(publicChatsMap[chat], msg)
	return nil
}

func (repo *Repository) GetPublicMessages(chat string) ([]models.Message, error) {
	messages, exists := publicChatsMap[chat]
	if !exists {
		return nil, errors.New("public chat not found")
	}
	return messages, nil
}

func (repo *Repository) GetAllPublicChats() ([]string, error) {
	return publicChatsSlice, nil
}

func (repo *Repository) AddPrivateMessage(receiver string, msg models.Message) error {
	// Check if chat exists for both users
	if privateChatsMap[receiver] == nil {
		privateChatsMap[receiver] = make(map[string]*models.PrivateChat)
	}
	if privateChatsMap[msg.Sender] == nil {
		privateChatsMap[msg.Sender] = make(map[string]*models.PrivateChat)
	}

	// Create new private chat if it doesn't exist
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

	// Add message to the private chat
	privateChatsMap[receiver][msg.Sender].Messages = append(privateChatsMap[receiver][msg.Sender].Messages, msg)
	return nil
}

func (repo *Repository) GetAllPrivateChats(user string) ([]string, error) {
	chats := make([]string, 0)
	for chat := range privateChatsMap[user] {
		chats = append(chats, chat)
	}
	return chats, nil
}

func (repo *Repository) GetPrivateMessages(userName, chatName string) ([]models.Message, error) {
	// Check if private chat exists for the user and chat
	if privateChatsMap[userName] == nil || privateChatsMap[userName][chatName] == nil {
		return nil, errors.New("private chat not found")
	}
	return privateChatsMap[userName][chatName].Messages, nil
}

func (repo *Repository) Close() error { return nil }
