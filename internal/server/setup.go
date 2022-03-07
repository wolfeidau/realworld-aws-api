package server

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/stores"
)

func Setup(cfg *flags.API, awscfg *aws.Config, e *echo.Echo) error {
	customerStore := stores.NewCustomers(awscfg, cfg)

	srv := NewCustomers(cfg, customerStore)

	// shut down all the default output of echo
	e.Logger.SetOutput(io.Discard)
	e.Logger.SetLevel(echolog.OFF)

	swagger, err := customersapi.GetSwagger()
	if err != nil {
		return errors.Wrap(err, "failed to loading swagger spec")
	}

	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				return nil // don't validate authentication options and enforce them as we use APIGW to do this
			},
		},
	}))

	customersapi.RegisterHandlers(e, srv)

	return nil
}
