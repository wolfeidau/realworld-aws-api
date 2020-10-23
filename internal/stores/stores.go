package stores

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/protobuf/proto"
	"github.com/wolfeidau/dynastore"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
)

type Record struct {
	Data    []byte
	Version int64
	ID      string
}

type Customers interface {
	GetCustomer(ctx context.Context, id string, into proto.Message) (int64, error)
	CreateCustomer(ctx context.Context, id, name string, obj proto.Message) (int64, error)
	ListCustomers(ctx context.Context, nextToken string, limit int) (string, []Record, error)
}

func NewCustomers(awsconfig *aws.Config, cfg *flags.API) Customers {

	session := dynastore.New(awsconfig)

	tbl := session.Table(cfg.CustomersTable)

	return &DDBCustomers{
		customerPart:      tbl.Partition(customersPartitionName),
		customerNamesPart: tbl.Partition(customerNamesPartitionName),
	}
}
