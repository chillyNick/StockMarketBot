CREATE TYPE user_state AS ENUM ('menu', 'add_stock', 'remove_stock');

CREATE TABLE "user"(
    id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    server_user_id INT NOT NULL,
    user_state user_state
);