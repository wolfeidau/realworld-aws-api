package commands

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

type NewCustomerCmd struct {
	Name        string   `kong:"required"`
	Labels      []string `kong:"required"`
	Description string
}

func (ccc *NewCustomerCmd) Run(ctx context.Context, cli *CLIContext) error {

	newCust := customersapi.NewCustomerJSONRequestBody{
		Description: nil,
		Labels:      ccc.Labels,
		Name:        ccc.Name,
	}

	if ccc.Description != "" {
		newCust.Description = &ccc.Description
	}

	res, err := cli.Customers.NewCustomerWithResponse(ctx, newCust)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusCreated {
		log.Ctx(ctx).Fatal().Str("status", res.Status()).Msg("request failed")
	}

	return cli.writeJson(res.JSON201)
}
