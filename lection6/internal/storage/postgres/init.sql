CREATE TABLE users
(
    id       serial PRIMARY KEY,
    name     text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE public_chats
(
    id   serial PRIMARY KEY,
    name text NOT NULL UNIQUE
);

CREATE TABLE public_messages
(
    id         serial PRIMARY KEY,
    chat_id    int REFERENCES public_chats (id) ON DELETE CASCADE,
    sender_id  int REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    content    text NOT NULL
);

CREATE TABLE private_chats
(
    id       serial PRIMARY KEY,
    user1_id int REFERENCES users (id),
    user2_id int REFERENCES users (id),
    UNIQUE (user1_id, user2_id)
);

CREATE TABLE private_messages
(
    id          serial PRIMARY KEY,
    chat_id     int REFERENCES private_chats (id) ON DELETE CASCADE,
    sender_id   int REFERENCES users (id),
    receiver_id int REFERENCES users (id),
    created_at  timestamp DEFAULT CURRENT_TIMESTAMP,
    content     text NOT NULL
);
