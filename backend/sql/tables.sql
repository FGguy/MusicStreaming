CREATE TABLE IF NOT EXISTS Users (
    name VARCHAR(30),
    password VARCHAR(50) NOT NULL,
    PRIMARY KEY(name)
);

CREATE INDEX IF NOT EXISTS idx_users_name
ON Users(name);