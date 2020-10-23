package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/cmd/customer-cli/apigw"
	"github.com/wolfeidau/realworld-aws-api/cmd/customer-cli/commands"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/httplog"
	"github.com/wolfeidau/realworld-aws-api/internal/logger"
)

var cfg struct {
	Version kong.VersionFlag
	Debug   bool
	URL     string `help:"The base URL for the API." kong:"required"`

	CreateCustomer commands.NewCustomerCmd   `cmd:"new-customer" help:"New Customer."`
	GetCustomer    commands.GetCustomerCmd   `cmd:"get-customer" help:"Read Customer."`
	ListCustomers  commands.ListCustomersCmd `cmd:"list-customers" help:"Read a list of Customers."`
}

func main() {
	cli := kong.Parse(&cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	log.Logger = logger.NewLogger()

	awscfg := new(aws.Config)

	httpClient := http.DefaultClient

	if cfg.Debug {
		httpClient.Transport = &httplog.Transport{}
	}

	client, err := customersapi.NewClientWithResponses(cfg.URL,
		customersapi.WithRequestEditorFn(apigw.RequestSigner(awscfg)), customersapi.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build client")
	}

	err = cli.Run(&commands.CLIContext{Customers: client, Debug: cfg.Debug, Writer: os.Stdout})
	cli.FatalIfErrorf(err)
}
