DROP TABLE sessions;

CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    expires_at TEXT,

    FOREIGN KEY (user_id) 
        REFERENCES users(user_id)
        ON UPDATE CASCADE 
        ON DELETE CASCADE
);
