
-- File created by: ./bin/db-tool new create_table_todo
BEGIN;
-- your migration here

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    title VARCHAR(255) NOT NULL,
    description VARCHAR(4096) NOT NULL
);

COMMIT;
