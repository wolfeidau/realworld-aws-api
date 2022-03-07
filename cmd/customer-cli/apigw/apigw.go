package apigw

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

func RequestSigner(awscfg *aws.Config) customersapi.RequestEditorFn {
	sess := session.Must(session.NewSession(awscfg))

	signer := v4.NewSigner(sess.Config.Credentials)

	return func(ctx context.Context, req *http.Request) error {
		log.Ctx(ctx).Info().Str("host", req.Host).Msg("signing request")

		body := bytes.NewReader([]byte{})

		if req.Body != nil {
			d, err := io.ReadAll(req.Body)
			if err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewReader(d))

			body = bytes.NewReader(d)
		}

		_, err := signer.Sign(req, body, "execute-api", aws.StringValue(sess.Config.Region), time.Now())
		if err != nil {
			return err
		}

		return nil
	}
}
