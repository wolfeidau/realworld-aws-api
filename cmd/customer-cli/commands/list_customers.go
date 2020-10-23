package commands

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

type ListCustomersCmd struct {
	NextToken string
	Limit     int `default:"100"`
}

func (ccc *ListCustomersCmd) Run(cli *CLIContext) error {

	log.Info().Msg("get a list of customers from the api")

	params := &customersapi.CustomersParams{
		MaxItems: &ccc.Limit,
	}

	if ccc.NextToken != "" {
		params.NextToken = &ccc.NextToken
	}

	res, err := cli.Customers.CustomersWithResponse(context.Background(), params)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list customers")
	}

	if res.StatusCode() != http.StatusOK {
		log.Fatal().Str("status", res.Status()).Msg("request failed")
	}

	return cli.writeJson(res.JSON200)
}
