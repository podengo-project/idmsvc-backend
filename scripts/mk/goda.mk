##
# Install goda locally to enable the usage
#
# goda generate a dependency diagram for golang code.
##

GODA ?= $(BIN)/goda

.PHONY: install-goda
install-goda: $(GODA)

GODA_VERSION ?= master

$(GODA):
	@{ \
	    export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
	    echo "Using GOPATH='$${GOPATH}'" ; \
	    [ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
	    export GOBIN="$(dir $(GODA))" ; \
	    echo "Installing 'goda' at '$(GODA)'" ; \
		pushd "$${GOPATH}" ; \
		git clone https://github.com/loov/goda && cd goda ; \
		go build -o "$(GODA)"; \
		popd ; \
	    find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
	    rm -rf "$${GOPATH}" ; \
	}
