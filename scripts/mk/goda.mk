##
# Install goda locally to enable the usage
#
# goda generate a dependency diagram for golang code.
##

GODA ?= $(BIN)/goda

.PHONY: install-goda
install-goda: $(GODA)

GODA_VERSION ?= v0.5.7

$(GODA):
	@{ \
	    export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
	    echo "Using GOPATH='$${GOPATH}'" ; \
	    [ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
	    export GOBIN="$(BIN)" ; \
	    echo "Installing 'goda' at '$(GODA)'" ; \
		go install "github.com/loov/goda@$(GODA_VERSION)"; \
	    find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
	    rm -rf "$${GOPATH}" ; \
	}
