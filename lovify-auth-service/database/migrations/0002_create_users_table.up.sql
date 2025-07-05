CREATE TABLE users(
    id bigserial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password varchar(64) NOT NULL,
);