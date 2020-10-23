package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
)

type CLIContext struct {
	Customers *customersapi.ClientWithResponses
	Writer    io.Writer
	Debug     bool
}

func (cc *CLIContext) writeJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(cc.Writer, string(data))
	if err != nil {
		return err
	}

	return nil
}
