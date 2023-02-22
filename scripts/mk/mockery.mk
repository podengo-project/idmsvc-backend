##
#
##

MOCKERY ?= $(BIN)/mockery

.PHONY: install-mockery
install-mockery: $(MOCKERY)

MOCKERY_VERSION ?= latest

$(MOCKERY):
	@{ \
	    export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
	    echo "Using GOPATH='$${GOPATH}'" ; \
	    [ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
	    export GOBIN="$(dir $(MOCKERY))" ; \
	    echo "Installing 'mockery' at '$(MOCKERY)'" ; \
	    go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION) ; \
	    find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
	    rm -rf "$${GOPATH}" ; \
	}
