
-- File created by: ./bin/db-tool new create_table_domains
BEGIN;
-- your migration here

DROP TABLE IF EXISTS ipa_locations;
DROP TABLE IF EXISTS ipa_servers;
DROP TABLE IF EXISTS ipa_certs;
DROP TABLE IF EXISTS ipas;
DROP TABLE IF EXISTS domains;

COMMIT;
