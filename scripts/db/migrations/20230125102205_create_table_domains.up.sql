
-- File created by: ./bin/db-tool new create_table_domains
BEGIN;
-- your migration here

-- See: https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-char-varchar-text/

-- NOTE https://samu.space/uuids-with-postgres-and-gorm/
--      thanks @anschnei
--      Consider to use UUID as the primary key
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    org_id      VARCHAR(255) NOT NULL,
    domain_uuid UUID UNIQUE NOT NULL,
    domain_name VARCHAR(253),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type INT NOT NULL,
    auto_enrollment_enabled BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS ipas (
    id INT UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    -- NOTE Keep in mind that gorm is making a logical delete,
    --      the row is not deleted from the database when
    --      using the normal operations.
    --      See: https://gorm.io/docs/delete.html
    deleted_at TIMESTAMP DEFAULT NULL,
    realm_name VARCHAR(253) NOT NULL,
    realm_domains TEXT NOT NULL,

    CONSTRAINT fk_ipas_id__domains_id
        FOREIGN KEY (id)
            REFERENCES domains(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ipa_certs (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    ipa_id INT,
    issuer TEXT NOT NULL,
    nickname VARCHAR(255) NOT NULL,
    not_after TIMESTAMP NOT NULL,
    not_before TIMESTAMP NOT NULL,
    serial_number VARCHAR(64) NOT NULL,
    subject TEXT NOT NULL,
    pem TEXT NOT NULL,

    CONSTRAINT fk_ipa_certs_ipa_id__ipas_id
        FOREIGN KEY (ipa_id)
            REFERENCES ipas(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ipa_servers (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    ipa_id INT,
    fqdn VARCHAR(253) NOT NULL,
    rhsm_id VARCHAR(64) DEFAULT NULL,
    location VARCHAR(63) DEFAULT NULL,
    ca_server BOOLEAN NOT NULL,
    hcc_enrollment_server BOOLEAN NOT NULL,
    hcc_update_server BOOLEAN NOT NULL,
    pk_init_server BOOLEAN NOT NULL,

    CONSTRAINT fk_ipa_servers_ipa_id__ipas_id
        FOREIGN KEY (ipa_id)
            REFERENCES ipas(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ipa_locations (
    id SERIAL UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    ipa_id INT,
    name VARCHAR(64) NOT NULL,
    description TEXT DEFAULT NULL,

    CONSTRAINT fk_ipa_locations_ipa_id__ipas_id
        FOREIGN KEY (ipa_id)
            REFERENCES ipas(id)
    ON DELETE CASCADE
);

COMMIT;
