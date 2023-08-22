##
# xrhidgen is a tool that helps to generate identities
# for the development on the console.dot platform for
# development and testing propose.
##
XRHIDGEN ?= $(BIN)/xrhidgen
XRHIDGEN_VERSION ?= latest

.PHONY: install-xrhidgen
install-xrhidgen: $(XRHIDGEN)

$(XRHIDGEN): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/subpop/xrhidgen/cmd/xrhidgen@$(XRHIDGEN_VERSION)"
