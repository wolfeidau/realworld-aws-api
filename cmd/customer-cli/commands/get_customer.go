package commands

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type GetCustomerCmd struct {
	ID string `kong:"required"`
}

func (ccc *GetCustomerCmd) Run(ctx context.Context, cli *CLIContext) error {
	log.Ctx(ctx).Info().Str("id", ccc.ID).Msg("get a customers from the api")

	res, err := cli.Customers.GetCustomerWithResponse(ctx, ccc.ID)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Fatal().Str("status", res.Status()).Msg("request failed")
	}

	return cli.writeJSON(res.JSON200)
}
