package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/wolfeidau/realworld-aws-api/internal/httplog"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

var cfg struct {
	Version kong.VersionFlag
	URL     string `help:"The base URL for the API." kong:"required"`

	CreateCustomer CreateCustomerCmd `cmd:"create-customer" help:"Create Customer."`
}

func main() {
	cli := kong.Parse(&cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Stack().Caller().Logger()

	awscfg := new(aws.Config)

	t := &httplog.Transport{}
	httpClient := &http.Client{Transport: t}

	client, err := customersapi.NewClientWithResponses(cfg.URL,
		customersapi.WithRequestEditorFn(requestSigner(awscfg)), customersapi.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build client")
	}

	err = cli.Run(&CLIContext{customers: client})
	cli.FatalIfErrorf(err)

	//log.Info().Msg("get a list of customers from the api")
	//
	//// to illustrate a basic client we can call to get a list of customers
	//res, err := client.CustomersWithResponse(context.Background(), &customersapi.CustomersParams{})
	//if err != nil {
	//	log.Fatal().Err(err).Msg("failed to list customers")
	//}
	//
	//if res.StatusCode() != http.StatusOK {
	//	log.Fatal().Str("status", res.Status()).Msg("request failed")
	//}
	//
	//log.Info().Fields(map[string]interface{}{
	//	"customerPage": res.JSON200,
	//}).Msg("customer list result")

}

func requestSigner(awscfg *aws.Config) customersapi.RequestEditorFn {

	sess := session.Must(session.NewSession(awscfg))

	signer := v4.NewSigner(sess.Config.Credentials)

	return func(ctx context.Context, req *http.Request) error {

		log.Info().Str("host", req.Host).Msg("signing request")

		body := bytes.NewReader([]byte{})

		if req.Body != nil {

			d, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(d))

			body = bytes.NewReader(d)
		}

		_, err := signer.Sign(req, body, "execute-api", aws.StringValue(sess.Config.Region), time.Now())
		if err != nil {
			return err
		}

		return nil
	}
}
