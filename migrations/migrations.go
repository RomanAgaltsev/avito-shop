// Package migrations implements database migrations.
package migrations

import "embed"

//go:embed "*.sql"
var Migrations embed.FS
