CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id          bigserial PRIMARY KEY,
    createdAt   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name        text                        NOT NULL,
    surname     text                        NOT NULL,
    username    CITEXT UNIQUE                      NOT NULL,
    email       text                        NOT NULL,
    password    BYTEA                       NOT NULL,
    activated   BOOL                        NOT NULL,
    version     INTEGER                     NOT NULL DEFAULT 1
    );

SELECT * FROM users;