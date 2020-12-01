package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/logger"
	"github.com/wolfeidau/realworld-aws-api/internal/migrate"
	"github.com/wolfeidau/realworld-aws-api/internal/server"
)

// extend the existing API flags with container settings
var cfg struct {
	flags.API

	Address          string `default:":3000" env:"ADDR"`
	DynamoDBEndpoint string `env:"DYNAMODB_ENDPOINT"`
}

func main() {
	// setup zerolog logger
	log.Logger = logger.NewLogger()

	kong.Parse(&cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	awscfg := new(aws.Config)

	if cfg.DynamoDBEndpoint != "" {
		log.Info().Str("endpoint", cfg.DynamoDBEndpoint).Msg("setting up local dynamodb")

		awscfg.Endpoint = aws.String(cfg.DynamoDBEndpoint)

		err := migrate.Table(awscfg, cfg.CustomersTable)
		if err != nil {
			log.Fatal().Err(err).Msg("database migration failed")
		}
	}

	e := echo.New()

	e.Use(logger.Middleware)

	err := server.Setup(&cfg.API, awscfg, e)
	if err != nil {
		log.Fatal().Err(err).Msg("server setup failed")
	}

	log.Info().Str("addr", cfg.Address).Msg("starting server")
	log.Fatal().Err(e.Start(cfg.Address)).Msg("start failed")
}
