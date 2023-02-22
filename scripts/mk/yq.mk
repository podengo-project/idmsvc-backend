# Rules to install yq tool

# https://github.com/mikefarah/yq#go-install
YQ ?= $(BIN)/yq
YQ_VERSION ?= v4.31.1

$(YQ):
	@{\
		export GOPATH="$(shell mktemp -d "$(PROJECT_DIR)/tmp.XXXXXXXX" 2>/dev/null)" ; \
		echo "Using GOPATH='$${GOPATH}'" ; \
		[ "$${GOPATH}" != "" ] || { echo "error:GOPATH is empty"; exit 1; } ; \
		export GOBIN="$(dir $(YQ))" ; \
		echo "Installing 'yq' at '$(YQ)'" ; \
		go install github.com/mikefarah/yq/v4@$(YQ_VERSION) ; \
		find "$${GOPATH}" -type d -exec chmod u+w {} \; ; \
		rm -rf "$${GOPATH}" ; \
	}

.PHONY: install-yq
install-yq: $(YQ)
