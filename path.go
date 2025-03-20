package mmdbgeo

import (
	"path/filepath"
	"runtime"
)

var (
	_, path, _, _ = runtime.Caller(0)
	rootDir       = filepath.Dir(path)
)
