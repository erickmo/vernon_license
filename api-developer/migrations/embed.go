// Package migrations menyimpan semua file SQL migration yang di-embed ke binary.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
