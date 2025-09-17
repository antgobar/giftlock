BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- users table
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username TEXT NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  created_at timestamptz DEFAULT now()
);

DROP TABLE IF EXISTS gifts CASCADE;
CREATE TABLE gifts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title text NOT NULL,
  description text,
  link text,
  price numeric(10,2),
  created_by uuid REFERENCES users(id) NOT NULL, -- who added it
  claimed_by uuid REFERENCES users(id) NULL, -- who claimed it (NULL if unclaimed)
  claimed_at timestamptz,
  created_at timestamptz DEFAULT now()
);

-- session table
CREATE TABLE IF NOT EXISTS sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) NOT NULL,
  token TEXT NOT NULL UNIQUE,
  created_at timestamptz DEFAULT now(),
  expires_at timestamptz NOT NULL
);

COMMIT;
