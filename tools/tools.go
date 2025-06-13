//go:build tools
// +build tools

package main

import (
	_ "github.com/achiku/planter"
	_ "github.com/atombender/go-jsonschema"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/loov/goda"
	_ "github.com/mikefarah/yq/v4"
	_ "github.com/spf13/cobra-cli"
	_ "github.com/subpop/xrhidgen/cmd/xrhidgen"
	_ "github.com/vektra/mockery/v2"
)
