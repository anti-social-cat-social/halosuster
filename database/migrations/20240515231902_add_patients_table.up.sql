CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE patient_gender
 AS ENUM (
'male',
'female'
);

CREATE TABLE IF NOT EXISTS patients (
identity_number VARCHAR UNIQUE PRIMARY KEY,
phone_number VARCHAR NOT NULL,
name VARCHAR NOT NULL,
birth_date date NOT NULL,
gender patient_gender NOT NULL,
identity_card_scan_img VARCHAR NOT NULL,
created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_patients_created_at ON patients(created_at);