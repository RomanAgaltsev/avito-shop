-- +goose Up
-- +goose StatementBegin
CREATE TABLE merch (
    id    SERIAL PRIMARY KEY,
    type  VARCHAR(20) UNIQUE NOT NULL,
    price INTEGER            NOT NULL DEFAULT 0,
);

CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(20) UNIQUE NOT NULL,
    password   VARCHAR(60)        NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT NOW()
);

CREATE TABLE balance
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL,
    coins    INTEGER            NOT NULL DEFAULT 0
);

CREATE TABLE inventory (
    id        SERIAL PRIMARY KEY,
    username  VARCHAR(20) NOT NULL,
    type      VARCHAR(20) NOT NULL,
    quantity  INTEGER     NOT NULL DEFAULT 0,
    bought_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

CREATE TABLE history (
    id        SERIAL PRIMARY KEY,
    username  VARCHAR(20) NOT NULL,
    from_user VARCHAR(20) NOT NULL,
    to_user   VARCHAR(20) NOT NULL,
    amount    INTEGER     NOT NULL DEFAULT 0,
    sent_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE history;
DROP TABLE inventory;
DROP TABLE balance;
DROP TABLE users;
DROP TABLE merch;
-- +goose StatementEnd