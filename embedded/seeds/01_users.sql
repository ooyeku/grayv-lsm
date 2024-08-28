-- Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email, password_hash) VALUES
('admin', 'admin@example.com', 'hashed_password_here'),
('user1', 'user1@example.com', 'hashed_password_here');

-- Down
DROP TABLE IF EXISTS users;