CREATE TABLE users (
    login       varchar NOT NULL UNIQUE,
    password       varchar NOT NULL
);

CREATE TABLE orders (
    number       varchar NOT NULL UNIQUE,
	uploaded_at timestamp,
    login       varchar REFERENCES users (login)
);

CREATE TABLE withdraws (
    order_num varchar unique,
	number numeric ,
	processed_at timestamp,
    login       varchar REFERENCES users (login)
);

CREATE TABLE internal_service (
    order_num varchar unique,
	status varchar ,
	accrual numeric
);

INSERT INTO users (login, password)
VALUES ('123', 'abracadabra');