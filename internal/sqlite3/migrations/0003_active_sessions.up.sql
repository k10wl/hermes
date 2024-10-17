DROP TABLE IF EXISTS active_sessions;

CREATE TABLE active_sessions (
    id INTEGER PRIMARY KEY,
    address TEXT NOT NULL,
    database_dns TEXT NOT NULL
);
