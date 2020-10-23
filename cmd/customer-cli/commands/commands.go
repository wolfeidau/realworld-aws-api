package commands

import (
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

type CLIContext struct {
	Customers *customersapi.ClientWithResponses
}
