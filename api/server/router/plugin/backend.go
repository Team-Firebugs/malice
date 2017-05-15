package plugin

import (
	"github.com/maliceio/malice/api/types"
	"github.com/maliceio/malice/api/types/filters"
)

// Backend for Plugin
type Backend interface {
	Disable(name string, config *types.Plugin) error
	Enable(name string, config *types.Plugin) error
	List(filters.Args) ([]types.Plugin, error)
	Inspect(name string) (*types.Plugin, error)
	Remove(name string, config *types.Plugin) error
}
