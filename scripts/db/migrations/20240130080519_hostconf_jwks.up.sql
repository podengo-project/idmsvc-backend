
-- File created by: bin/db-tool new hostconf_jwks
BEGIN;

CREATE TABLE IF NOT EXISTS hostconf_jwks (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    key_id VARCHAR(16) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    public_jwk TEXT NOT NULL,
    encryption_id VARCHAR(16) NOT NULL,
    encrypted_jwk BYTEA
);

COMMIT;
