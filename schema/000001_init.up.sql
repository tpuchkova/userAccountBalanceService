CREATE TABLE transactions
(
  id serial not null unique,
  state varchar(255) not null,
  amount numeric(25,2) not null,
  transaction_id varchar(255) not null unique,
  canceled boolean NOT NULL DEFAULT false,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_balance
(
    id serial not null unique,
    balance numeric(25,2) not null
);

INSERT INTO user_balance (balance)
VALUES (10);