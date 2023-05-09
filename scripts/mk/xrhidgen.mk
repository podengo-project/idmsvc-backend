##
# xrhidgen is a tool that helps to generate identities
# for the development on the console.dot platform for
# development and testing propose.
##
XRHIDGEN ?= $(BIN)/xrhidgen
XRHIDGEN_VERSION ?= latest

.PHONY: install-xrhidgen
install-xrhidgen: $(XRHIDGEN)

$(XRHIDGEN):
	@{ \
	    export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
	    echo "Using GOPATH='$${GOPATH}'" ; \
	    [ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
	    export GOBIN="$(dir $(XRHIDGEN))" ; \
	    echo "Installing 'xrhidgen' at '$(XRHIDGEN)'" ; \
	    go install github.com/subpop/xrhidgen/cmd/xrhidgen@$(XRHIDGEN_VERSION) ; \
	    find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
	    rm -rf "$${GOPATH}" ; \
	}
