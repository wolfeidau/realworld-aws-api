package main

import (
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/server"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	srv := server.NewCustomers(cfg)

	// build a list of fields to include in all events
	flds := lmw.FieldMap{"commit": app.Commit, "buildDate": app.BuildDate}

	e := echo.New()

	// shut down all the default output of echo
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	swagger, err := customersapi.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to loading swagger spec")
	}

	customersapi.RegisterHandlers(e, srv, middleware.OapiRequestValidator(swagger))

	gw := gateway.NewGateway(e)

	ch := lmw.New(zlog.New(zlog.Fields(flds))).Then(gw)

	lambda.StartHandler(ch)
}
