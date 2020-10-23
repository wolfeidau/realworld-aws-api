package commands

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type GetCustomerCmd struct {
	ID string `kong:"required"`
}

func (ccc *GetCustomerCmd) Run(cli *CLIContext) error {

	res, err := cli.Customers.GetCustomerWithResponse(context.Background(), ccc.ID)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		log.Fatal().Str("status", res.Status()).Msg("request failed")
	}

	log.Info().Fields(map[string]interface{}{
		"customer": res.JSON200,
	}).Msg("customer read")

	return nil
}
