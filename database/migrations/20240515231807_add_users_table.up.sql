CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role
 AS ENUM (
'it',
'nurse'
);

CREATE TABLE IF NOT EXISTS users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
nip VARCHAR UNIQUE NOT NULL,
role user_role NOT NULL,
name VARCHAR NOT NULL,
password VARCHAR,
identity_card_scan_img VARCHAR,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);