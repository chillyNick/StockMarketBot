CREATE TYPE user_state AS ENUM ('menu', 'add_stock', 'remove_stock, diff, add_notification, remove_notification');

CREATE TABLE "user"(
    id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    server_user_id INT NOT NULL,
    state user_state
);