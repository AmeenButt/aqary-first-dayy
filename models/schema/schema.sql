CREATE TABLE IF NOT EXISTS users (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  email text,
  password text,
  profile_picture text,
  otp int,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_wallet (
  id   BIGSERIAL PRIMARY KEY,
  user_id int,
  amount float default 0.0,
  FOREIGN KEY (user_id) REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_transactions (
  id   BIGSERIAL PRIMARY KEY,
  action text not null,
  user_wallet_id int,
  transaction_amount float,
  FOREIGN KEY (user_wallet_id) REFERENCES user_wallet(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS properties (
  id   BIGSERIAL PRIMARY KEY,
  sizeInSqFeet int,
  location text,
  images text[],
  demand text,
  status text DEFAULT 'pending',
  user_id int,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
