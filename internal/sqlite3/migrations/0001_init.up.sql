-- Table for storing chats
CREATE TABLE chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Primary key
    name TEXT DEFAULT "unnamed",
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Table for storing messages
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Primary key
    chat_id INTEGER NOT NULL, -- Foreign key to chats table
    content TEXT NOT NULL,
    role_id INTEGER NOT NULL, -- Foreign key to roles table
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- Table for storing roles
CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Auto-increment primary key
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE web_settings (
    dark_mode BOOLEAN NOT NULL,
    initted BOOLEAN NOT NULL
);

-- Inserting values into roles table
INSERT INTO roles (name) VALUES ('user'), ('assistant'), ('system');

INSERT INTO web_settings (dark_mode, initted) VALUES (false, false);
