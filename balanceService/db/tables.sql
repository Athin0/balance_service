CREATE SCHEMA IF NOT EXISTS balance;

CREATE TABLE IF NOT EXISTS balance.balance
(
    user_id bigserial PRIMARY KEY,
    value decimal(10, 2) NOT NULL CHECK (value >= 0) DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS balance.history
(
    id         bigserial PRIMARY KEY,
    user_id bigint NOT NULL ,
    service_id bigint NOT NULL ,
    order_id bigint NOT NULL ,
    value decimal(10, 2) NOT NULL CHECK (value >= 0),
    occurred_at timestamptz NOT NULL,
    description text
);

CREATE TABLE IF NOT EXISTS balance.reserved
(
    id          bigserial PRIMARY KEY,
    user_id bigint NOT NULL ,
    service_id bigint NOT NULL ,
    order_id bigint NOT NULL ,
    value decimal(10, 2) NOT NULL CHECK (value >= 0)
)