package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/server"
	"github.com/wolfeidau/realworld-aws-api/internal/stores"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	awscfg := new(aws.Config)

	customerStore := stores.NewCustomers(awscfg, cfg)

	srv := server.NewCustomers(cfg, customerStore)

	// build a list of fields to include in all events
	flds := lmw.FieldMap{"commit": app.Commit, "buildDate": app.BuildDate, "stage": cfg.Stage, "branch": cfg.Branch}

	e := echo.New()

	// shut down all the default output of echo
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	swagger, err := customersapi.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to loading swagger spec")
	}

	customersapi.RegisterHandlers(e, srv, middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				return nil // don't validate authentication options and enforce them as we use APIGW to do this
			},
		},
	}))

	gw := gateway.NewGateway(e)

	ch := lmw.New(
		zlog.New(zlog.Fields(flds)), // build a logger and inject it into the context
	)

	if cfg.Stage == "dev" {
		ch.Use(raw.New(raw.Fields(flds))) // raw event logger used during development
	}

	h := ch.Then(gw)

	lambda.StartHandler(h)
}
