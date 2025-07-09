CREATE TABLE tokens(
    id bigserial PRIMARY KEY,
    token text NOT NULL,
    type token_type NOT NULL,
    expiration_date varchar(64),
    email text references users(email)
);