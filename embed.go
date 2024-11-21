package nginfier

import (
	"embed"
)

//go:embed all:templates
var templates embed.FS

func Templates() embed.FS {
	return templates
}
