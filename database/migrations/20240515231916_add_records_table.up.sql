CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS records (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_by UUID NOT NULL REFERENCES users(id),
identity_number VARCHAR NOT NULL REFERENCES patients(identity_number),
symptoms VARCHAR NOT NULL,
medications VARCHAR NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_records_created_at ON records(created_at);