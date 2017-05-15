package image

import "github.com/maliceio/malice/api/types"

// Backend is the methods that need to be implemented to provide
// search specific functionality
type Backend interface {
	Scans(filter string) ([]*types.Scan, []string, error)
	Pipe(id, path string) ([]*types.Scan, []string, error)
}
