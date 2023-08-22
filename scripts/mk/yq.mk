# Rules to install yq tool

# https://github.com/mikefarah/yq#go-install
YQ ?= $(BIN)/yq
YQ_VERSION ?= v4.31.1

$(YQ): $(BIN)
	GOBIN="$(dir $(CURDIR)/$@)" go install "github.com/mikefarah/yq/v4@$(YQ_VERSION)"

.PHONY: install-yq
install-yq: $(YQ)
