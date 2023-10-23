##
# Rules to browse the openapi specification by using redocly tool
#
# See: https://hub.docker.com/r/redocly/redoc/
##

OPENAPI_FILE ?= $(PROJECT_DIR)/api/public.openapi.yaml

.PHONY: open-api-doc
open-api-doc:  ## Open OPENAPI_FILE to browse the documentation (by default the public openapi specification)
	$(CONTAINER_ENGINE) run -it --name open-api-doc -d --rm -p 8080:80 -v "$(OPENAPI_FILE):/usr/share/nginx/html/swagger.yaml:ro,z" -e SPEC_URL=swagger.yaml docker.io/redocly/redoc
	xdg-open http://localhost:8080

.PHONY: stop-api-doc
stop-api-doc:  ## Stop open-api-doc container
	$(CONTAINER_ENGINE) stop open-api-doc

