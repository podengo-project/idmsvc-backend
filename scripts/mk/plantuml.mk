##
# Rules related with the generation of plantuml diagrams.
#
# NOTE: Keep in mind that they don't need to be added to the
#       repository as it can be seen at the link below:
#       https://blog.anoff.io/2018-07-31-diagrams-with-plantuml/
#
# NOTE: You need to install plantuml by hand:
#       In Fedora systems you can do that by:
#       # dnf install -y plantuml
#
# PLANTER_NO_GENERATE=y will avoid to run planter to generate docs/db-model.puml
##

PLANTUML ?= $(shell command -v plantuml 2>/dev/null)
PLANTUML ?= false

PLANTUML_SOURCES ?= $(patsubst docs/%.puml,docs/%.svg,$(wildcard docs/*.puml)) $(patsubst docs/sequence/%.puml,docs/sequence/%.svg,$(wildcard docs/sequence/*.puml))
PLANTER_NO_GENERATE ?= n
.PHONY: generate-diagrams
generate-diagrams: $(PLANTUML_SOURCES)  ## Generate diagrams (PLANTER_NO_GENERATE=y to don't generate docs/db-model.puml)
ifneq (y,$(PLANTER_NO_GENERATE))
	$(MAKE) docs/db-model.svg
endif

docs/db-model.puml: $(PLANTER) .compose-wait-db scripts/db/migrations/*.up.sql
	$(PLANTER) -o $@ postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST)/$(DATABASE_NAME)?sslmode=disable

# General rule to generate a diagram in SVG format for
# each .puml file found at docs/ directory
docs/%.svg: docs/%.puml
	$(PLANTUML) -tsvg $<
docs/sequence/%.svg: docs/sequence/%.puml
	$(PLANTUML) -tsvg $<
