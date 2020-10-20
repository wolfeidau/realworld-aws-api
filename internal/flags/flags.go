package flags

import "github.com/alecthomas/kong"

// API api related flags passing in env variables
type API struct {
	Version kong.VersionFlag
}
