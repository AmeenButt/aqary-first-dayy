CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  email text,
  password text,
  author_status boolean default false,
  author_id int,
  FOREIGN KEY (author_id) REFERENCES authors(id)
);

CREATE TABLE user_wallet (
  id   BIGSERIAL PRIMARY KEY,
  user_id int,
  amount text,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE user_transactions (
  id   BIGSERIAL PRIMARY KEY,
  user_wallet_id int,
  transaction_amount text,
  FOREIGN KEY (user_wallet_id) REFERENCES user_wallet(id)
);