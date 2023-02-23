
-- File created by: ./bin/db-tool new create_table_domains
BEGIN;
-- your migration here

-- NOTE https://samu.space/uuids-with-postgres-and-gorm/
--      thanks @anschnei
--      Consider to use UUID as the primary key
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    domain_uuid UUID UNIQUE NOT NULL,
    domain_name VARCHAR(253) NOT NULL,
    domain_type INT NOT NULL,
    auto_enrollment_enabled BOOLEAN NOT NULL,
);

CREATE TABLE IF NOT EXISTS ipas (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    -- NOTE Keep in mind that gorm is making a logical delete,
    --      the row is not deleted from the database when
    --      using the normal operations.
    --      See: https://gorm.io/docs/delete.html
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,

    domain_id INT,
    realm_name VARCHAR(253) NOT NULL,
    -- See: https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-char-varchar-text/
    ca_list TEXT NOT NULL,
    server_list TEXT NOT NULL
);

ALTER TABLE ipas
ADD CONSTRAINT fk_domain
FOREIGN KEY (domain_id)
REFERENCES domains(id)
ON DELETE SET NULL;

COMMIT;
