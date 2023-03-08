##
# The target for this file is just to enumerate which partial
# makefile we want to use to compose our final Makefile.
#
# Unless you are not using conditional assignment within
# the different variable files, this would be the priority:
# - The values indicated at 'configs/config.yaml' file.
# - The values indicated at 'mk/variables.mk' file. This
#   file is included into the repository and define the
#   default values for the variables, if not assigned yet.
# - The 'mk/meta-*.mk' files just contain the comment to
#   print out the group text for the help content. They
#   are into independent files, because the order they
#   appear into this include file matters, and provide
#   the flexibility to print out the group text exactly
#   where we want kust changing the order into this file.
#
# This file set the 'help' rule as the default one when
# no arguments are indicated.
##
include scripts/mk/projectdir.mk
-include secrets/private.mk
include scripts/mk/variables.mk

include scripts/mk/help.mk
include scripts/mk/meta-general.mk
include scripts/mk/gojsonschema.mk
include scripts/mk/yq.mk
# include scripts/mk/gen-api.mk
include scripts/mk/go-rules.mk
include scripts/mk/api-doc.mk
include scripts/mk/db.mk
include scripts/mk/printvars.mk
include scripts/mk/plantuml.mk
# include scripts/mk/swag.mk
include scripts/mk/lint.mk
include scripts/mk/meta-docker.mk
include scripts/mk/docker.mk
include scripts/mk/meta-compose.mk
include scripts/mk/compose.mk
include scripts/mk/meta-kafka.mk
include scripts/mk/kafka.mk
# include scripts/mk/meta-repos.mk
# include scripts/mk/repos.mk
include scripts/mk/meta-ephemeral.mk
include scripts/mk/ephemeral.mk

include scripts/mk/mockery.mk
include scripts/mk/oapi-codegen.mk
include scripts/mk/goda.mk
include scripts/mk/meta-prometheus.mk
include scripts/mk/prometheus.mk
