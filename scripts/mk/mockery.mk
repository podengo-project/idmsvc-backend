##
#
##

MOCKERY ?= $(BIN)/mockery

.PHONY: install-mockery
install-mockery: $(MOCKERY)

MOCKERY_VERSION ?= v2.16.0

$(MOCKERY): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/vektra/mockery/v2@$(MOCKERY_VERSION)"
