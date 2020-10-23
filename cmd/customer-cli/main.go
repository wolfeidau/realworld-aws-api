package main

import (
	"fmt"
	"net/http"

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

	CreateCustomer commands.NewCustomerCmd `cmd:"new-customer" help:"New Customer."`
	GetCustomer    commands.GetCustomerCmd `cmd:"get-customer" help:"Read Customer."`
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

	err = cli.Run(&commands.CLIContext{Customers: client})
	cli.FatalIfErrorf(err)

	//log.Info().Msg("get a list of customers from the api")
	//
	//// to illustrate a basic client we can call to get a list of customers
	//res, err := client.CustomersWithResponse(context.Background(), &customersapi.CustomersParams{})
	//if err != nil {
	//	log.Fatal().Err(err).Msg("failed to list customers")
	//}
	//
	//if res.StatusCode() != http.StatusOK {
	//	log.Fatal().Str("status", res.Status()).Msg("request failed")
	//}
	//
	//log.Info().Fields(map[string]interface{}{
	//	"customerPage": res.JSON200,
	//}).Msg("customer list result")

}
