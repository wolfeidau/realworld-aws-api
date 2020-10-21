package stores

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/protobuf/proto"
	"github.com/wolfeidau/dynastore"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
)

type Customers interface {
	GetCustomer(ctx context.Context, id string, into proto.Message) (int64, error)
	CreateCustomer(ctx context.Context, id string, obj proto.Message) (int64, error)
}

func NewCustomers(awsconfig *aws.Config, cfg *flags.API) Customers {

	session := dynastore.New(awsconfig)

	return &DDBCustomers{
		session.Table(cfg.CustomersTable).Partition("customers"),
	}
}
