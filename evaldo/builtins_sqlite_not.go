//go:build no_sqlite
// +build no_sqlite

package evaldo

import (
	"github.com/refaktor/rye/env"
)

var Builtins_sqlite = map[string]*env.Builtin{}
