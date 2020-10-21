package flags

import "github.com/alecthomas/kong"

// API api related flags passing in env variables
type API struct {
	Version        kong.VersionFlag
	AppName        string `help:"Stage the name of the service." env:"APP_NAME"`
	Stage          string `help:"Stage the software is deployed." env:"STAGE"`
	Branch         string `help:"Branch used to build software." env:"BRANCH"`
	CustomersTable string `help:"Name of the dynamodb to use for storing customers." env:"CUSTOMERS_TABLE"`
}

// Client client flags
type Client struct {
	Version kong.VersionFlag
	URL     string `help:"The base URL for the API." kong:"required"`
}
