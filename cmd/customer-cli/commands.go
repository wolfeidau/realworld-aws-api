package main

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

type CLIContext struct {
	customers *customersapi.ClientWithResponses
}

type CreateCustomerCmd struct {
	Name        string   `kong:"required"`
	Labels      []string `kong:"required"`
	Description string
}

func (ccc *CreateCustomerCmd) Run(cli *CLIContext) error {

	newCust := customersapi.NewCustomerJSONRequestBody{
		Description: nil,
		Labels:      ccc.Labels,
		Name:        ccc.Name,
	}

	if ccc.Description != "" {
		newCust.Description = &ccc.Description
	}

	res, err := cli.customers.NewCustomerWithResponse(context.Background(), newCust)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusCreated {
		log.Fatal().Str("status", res.Status()).Msg("request failed")
	}

	log.Info().Fields(map[string]interface{}{
		"customer": res.JSON201,
	}).Msg("customer created")

	return nil
}
