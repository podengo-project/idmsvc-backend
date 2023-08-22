##
# Install goda locally to enable the usage
#
# goda generate a dependency diagram for golang code.
##

GODA ?= $(BIN)/goda

.PHONY: install-goda
install-goda: $(GODA)

GODA_VERSION ?= v0.5.7

$(GODA): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/loov/goda@$(GODA_VERSION)"
