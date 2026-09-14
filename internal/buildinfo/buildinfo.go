// Package buildinfo exposes metadata embedded into the executable at build time.
package buildinfo

import "fmt"

// These defaults make local development builds useful. Release builds can
// replace them with go build -ldflags "-X import/path.Name=value".
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// String returns a stable, human-readable representation of the build.
func String() string {
	return fmt.Sprintf("version=%s commit=%s built=%s", Version, Commit, BuildTime)
}
