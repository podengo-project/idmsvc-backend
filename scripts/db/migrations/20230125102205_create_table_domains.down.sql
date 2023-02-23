
-- File created by: ./bin/db-tool new create_table_domains
BEGIN;
-- your migration here

DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS ipas;

COMMIT;
