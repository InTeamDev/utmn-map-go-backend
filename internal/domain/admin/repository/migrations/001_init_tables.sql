CREATE TABLE administrators (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    full_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    salt VARCHAR(100), 
    password_reset_token VARCHAR(255) NULL,
    token_expires_at TIMESTAMP NULL
);

-- Создание индексов для ускорения поиска
CREATE INDEX idx_administrators_username ON administrators(username);
CREATE INDEX idx_administrators_email ON administrators(email);

---- create above / drop below ----

DROP TABLE administrators;
