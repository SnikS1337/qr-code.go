package site_files

import "embed"

// IndexHTML using embed filesystem in handlers.go
//
//go:embed index.html
var IndexHTML embed.FS
