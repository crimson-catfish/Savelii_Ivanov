INSERT INTO users (name, password)
VALUES ('bill', 'bill_rules');
INSERT INTO users (name, password)
VALUES ('uma_thruman', 'bill_sucks');

INSERT INTO public_chats (name)
VALUES ('greatest-chat-of-all-times');
INSERT INTO public_chats (name)
VALUES ('ok-chat');
INSERT INTO public_chats (name)
VALUES ('meh-chat');

INSERT INTO public_messages (chat_id, sender_id, content)
VALUES (1, 1, 'welcome to the greatest chat of all times!!!');

INSERT INTO private_chats (user1_id, user2_id)
VALUES (1, 2);

INSERT INTO private_messages (chat_id, sender_id, receiver_id, content)
VALUES (1, 2, 1, 'ima kill ya');
