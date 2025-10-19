// Package template is used to embed the default templates into the binary.
package template

import (
	"embed"
)

var (
	// DefaultTemplates is the default template file system.
	//go:embed all:*.tmpl
	DefaultTemplates embed.FS
)
