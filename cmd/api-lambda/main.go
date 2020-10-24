package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/labstack/echo/v4"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/server"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	awscfg := new(aws.Config)

	e := echo.New()
	err := server.Setup(cfg, awscfg, e)
	if err != nil {
		log.Fatal().Err(err).Msg("server setup failed")
	}

	gw := gateway.NewGateway(e)

	// build a list of fields to include in all events
	flds := lmw.FieldMap{"commit": app.Commit, "buildDate": app.BuildDate, "stage": cfg.Stage, "branch": cfg.Branch}

	ch := lmw.New(
		zlog.New(zlog.Fields(flds)), // build a logger and inject it into the context
	)

	if cfg.Stage == "dev" {
		ch.Use(raw.New(raw.Fields(flds))) // raw event logger used during development
	}

	h := ch.Then(gw)

	lambda.StartHandler(h)
}
