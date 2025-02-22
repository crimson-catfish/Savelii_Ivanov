package postgres

import (
	"context"
	"log"

	"entrance/lection6/internal/models"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewPostgresRepository() *Repository {
	dsn := "postgres://username:password@localhost:5442/chat?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) UserExists(name string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM users WHERE name = $1);`
	var exists bool
	if err := r.conn.QueryRow(context.Background(), query, name).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
func (r *Repository) AddUser(credentials models.Credentials) error {
	const query = `INSERT INTO users (name, password) VALUES ($1, $2);`
	_, err := r.conn.Exec(context.Background(), query, credentials.Name, credentials.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPassword(name string) (string, error) {
	const query = `SELECT password FROM users WHERE name = $1;`
	var password string
	if err := r.conn.QueryRow(context.Background(), query, name).Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}
func (r *Repository) GetAllPublicChats() ([]string, error) {
	const query = `SELECT name FROM public_chats;`
	rows, err := r.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publicChats := make([]string, 0)
	for rows.Next() {
		var chatName string
		if err := rows.Scan(&chatName); err != nil {
			return nil, err
		}
		publicChats = append(publicChats, chatName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return publicChats, nil
}

func (r *Repository) GetPublicMessages(chat string) ([]models.Message, error) {
	query := `
		SELECT 
			u.name AS sender,
			pm.created_at AS time,
			pm.content
		FROM 
			public_chats pc
		JOIN 
			public_messages pm ON pc.id = pm.chat_id
		JOIN 
			users u ON pm.sender_id = u.id
		WHERE 
			pc.name = $1
		ORDER BY 
			pm.created_at;
	`

	rows, err := r.conn.Query(context.Background(), query, chat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Sender, &message.Time, &message.Content); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) AddPublicMessage(chat string, msg models.Message) error {
	const query = `
		INSERT INTO public_messages (chat_id, sender_id, content)
		SELECT pc.id, u.id, $1
		FROM public_chats pc
		JOIN users u ON u.name = $2
		WHERE pc.name = $3;
	`
	_, err := r.conn.Exec(context.Background(), query, msg.Content, msg.Sender, chat)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllPrivateChats(user string) ([]string, error) {
	const query = `
		SELECT 
			CASE 
				WHEN pc.user1_id = u.id THEN u2.name
				WHEN pc.user2_id = u.id THEN u1.name
			END AS chat_name
		FROM private_chats pc
		JOIN users u1 ON pc.user1_id = u1.id
		JOIN users u2 ON pc.user2_id = u2.id
		JOIN users u ON u.name = $1
	`
	rows, err := r.conn.Query(context.Background(), query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []string
	for rows.Next() {
		var chatName *string
		if err := rows.Scan(&chatName); err != nil {
			return nil, err
		}
		if chatName != nil {
			chats = append(chats, *chatName)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *Repository) GetPrivateMessages(userName, chatName string) ([]models.Message, error) {
	const query = `
		SELECT 
			u.name AS sender,
			pm.created_at AS time,
			pm.content
		FROM private_chats pc
		JOIN private_messages pm ON pc.id = pm.chat_id
		JOIN users u ON pm.sender_id = u.id
		JOIN users u1 ON u1.name = $1
		JOIN users u2 ON u2.name = $2
		WHERE 
			(pc.user1_id = u1.id AND pc.user2_id = u2.id) OR
			(pc.user1_id = u2.id AND pc.user2_id = u1.id)
		ORDER BY pm.created_at;
	`

	rows, err := r.conn.Query(context.Background(), query, userName, chatName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Sender, &message.Time, &message.Content); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) AddPrivateMessage(receiver string, msg models.Message) error {
	// Find the private chat ID between the sender and receiver
	const findChatQuery = `
        SELECT pc.id
        FROM private_chats pc
        JOIN users u1 ON pc.user1_id = u1.id
        JOIN users u2 ON pc.user2_id = u2.id
        WHERE (u1.name = $1 AND u2.name = $2) OR (u1.name = $2 AND u2.name = $1);
    `

	var chatID int
	err := r.conn.QueryRow(context.Background(), findChatQuery, msg.Sender, receiver).Scan(&chatID)
	if err != nil {
		return err
	}

	// Find the sender's ID
	const findSenderIDQuery = `SELECT id FROM users WHERE name = $1;`
	var senderID int
	err = r.conn.QueryRow(context.Background(), findSenderIDQuery, msg.Sender).Scan(&senderID)
	if err != nil {
		return err
	}

	// Find the receiver's ID
	const findReceiverIDQuery = `SELECT id FROM users WHERE name = $1;`
	var receiverID int
	err = r.conn.QueryRow(context.Background(), findReceiverIDQuery, receiver).Scan(&receiverID)
	if err != nil {
		return err
	}

	// Insert the message into the private_messages table
	const insertMessageQuery = `
        INSERT INTO private_messages (chat_id, sender_id, receiver_id, content)
        VALUES ($1, $2, $3, $4);
    `
	_, err = r.conn.Exec(context.Background(), insertMessageQuery, chatID, senderID, receiverID, msg.Content)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Close() error {
	err := r.conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}
