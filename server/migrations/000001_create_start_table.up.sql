-- Создание начальных таблиц приложения


-- TABLE users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    phone_number TEXT,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    lastactive_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);
CREATE INDEX idx_users_email ON users(email);


-- TABLE bytes
CREATE TABLE bytes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    sent_size BIGINT,
    received_size BIGINT,
    type VARCHAR(15) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);


-- TABLE texts
CREATE TABLE texts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type VARCHAR(30) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);

--  name TEXT NOT NULL,
--  должно быть уникально в рамках
--   user_email INTEGER REFERENCES users(id) ON DELETE CASCADE