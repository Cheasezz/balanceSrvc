CREATE TABLE IF NOT EXISTS users (
  id        UUID            PRIMARY KEY,
  balance   BIGINT          NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction_types (
  id        SMALLSERIAL  PRIMARY KEY,
  code      TEXT         UNIQUE         NOT NULL,
  name      TEXT         NOT NULL,
  category  TEXT         NOT NULL       CHECK(category = lower(category))
);

INSERT INTO transaction_types (code, name, category) VALUES 
  ('deposit', 'Deposit', 'system'),
  ('withdrawal', 'Withdrawal', 'system'),
  ('reward', 'Reward', 'system'),
  ('transfer', 'User transfer', 'user')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE IF NOT EXISTS transactions (
  id            UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_id     UUID,
  resipient_id  UUID,
  type_id       SMALLINT        NOT NULL,
  amount        BIGINT          NOT NULL,
  created_at    TIMESTAMP       NOT NULL    DEFAULT now(),

  CONSTRAINT sender_fk            FOREIGN KEY (sender_id)      REFERENCES users               ON DELETE SET NULL,
  CONSTRAINT resipient_fk         FOREIGN KEY (resipient_id)   REFERENCES users               ON DELETE SET NULL,
  CONSTRAINT transaction_types_fk FOREIGN KEY (type_id)        REFERENCES transaction_types
);
