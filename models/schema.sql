CREATE TABLE users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  email text,
  password text
);

CREATE TABLE user_wallet (
  id   BIGSERIAL PRIMARY KEY,
  user_id int,
  amount float default 0.0,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_transactions (
  id   BIGSERIAL PRIMARY KEY,
  action text not null,
  user_wallet_id int,
  transaction_amount float,
  FOREIGN KEY (user_wallet_id) REFERENCES user_wallet(id)
);
