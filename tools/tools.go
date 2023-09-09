//go:build tools
// +build tools

package main

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/loov/goda"
	_ "github.com/xeipuuv/gojsonschema"
	_ "github.com/vektra/mockery/v2"
	_ "github.com/atombender/go-jsonschema/cmd/gojsonschema"
	_ "github.com/achiku/planter"
	_ "github.com/subpop/xrhidgen/cmd/xrhidgen"
	_ "github.com/mikefarah/yq/v4"
)
