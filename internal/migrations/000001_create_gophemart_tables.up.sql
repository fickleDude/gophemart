CREATE TABLE IF NOT EXISTS users (
    login       varchar NOT NULL UNIQUE,
    password       varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    number       varchar NOT NULL UNIQUE,
	uploaded_at timestamp,
    login       varchar REFERENCES users (login)
);

CREATE TABLE IF NOT EXISTS withdraws (
    order_num varchar unique,
	sum numeric ,
	processed_at timestamp,
    login       varchar REFERENCES users (login)
);