CREATE DATABASE orders;
CREATE DATABASE billing;
CREATE DATABASE warehouse;

\c orders;

CREATE TABLE orders (
    id            SERIAL PRIMARY KEY,
    item_id       INTEGER NOT NULL,
    user_id       INTEGER NOT NULL,
    payment_id    INTEGER,
    status        TEXT NOT NULL CHECK(status IN ('PENDING', 'FAILED', 'CREATED')),
    created_at    TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP at time zone 'utc')
);

\c billing;

CREATE TABLE payments (
    id            SERIAL PRIMARY KEY,
    order_id      INTEGER NOT NULL,
    from_user     INTEGER NOT NULL,
    to_user       INTEGER NOT NULL,
    amount        INTEGER NOT NULL
);

\c warehouse;

CREATE TABLE items (
    id            SERIAL PRIMARY KEY,
    seller_id     INTEGER NOT NULL,
    price         INTEGER NOT NULL,  -- better use MONEY type
    quantity      INTEGER NOT NULL
);

ALTER ROLE postgres SET client_encoding TO 'utf8';
ALTER ROLE postgres SET default_transaction_isolation TO 'read committed';
ALTER ROLE postgres SET timezone TO 'UTC';

GRANT ALL PRIVILEGES ON DATABASE orders TO postgres;
GRANT ALL PRIVILEGES ON DATABASE billing TO postgres;
GRANT ALL PRIVILEGES ON DATABASE warehouse TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
