CREATE TABLE "user"(
   id serial PRIMARY KEY
);

CREATE TABLE "stock" (
    id serial PRIMARY KEY,
    name VARCHAR NOT NULL,
    user_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES "user"(id),
    UNIQUE (name, user_id)
);

CREATE TYPE notification_type AS ENUM ('up', 'down');

CREATE TABLE "notification" (
     id serial PRIMARY KEY,
     stock_name VARCHAR NOT NULL,
     user_id INTEGER NOT NULL,
     threshold DOUBLE PRECISION NOT NULL,
     type notification_type,
     CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES "user"(id)
)