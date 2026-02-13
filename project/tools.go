//go:build tools
// +build tools

package tools

import (
	_ "github.com/pressly/goose/v3"
	_ "github.com/sqlc-dev/sqlc"

	// format
	_ "golang.org/x/tools/go/packages"
)
