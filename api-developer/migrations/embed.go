package migrations

import "embed"

// FS berisi semua file SQL migrasi, di-embed ke dalam binary.
//
//go:embed *.sql
var FS embed.FS
