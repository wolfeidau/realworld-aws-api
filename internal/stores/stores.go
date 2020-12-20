package stores

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/wolfeidau/dynastore"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"google.golang.org/protobuf/proto"
)

// Record used to return raw record data
type Record struct {
	Data    []byte
	Version int64
	ID      string
}

// Customers store used to manage customers
type Customers interface {
	GetCustomer(ctx context.Context, id string, into proto.Message) (int64, error)
	CreateCustomer(ctx context.Context, id, name string, obj proto.Message) (int64, error)
	ListCustomers(ctx context.Context, nextToken string, limit int) (string, []Record, error)
}

// NewCustomers creates a new customer store
func NewCustomers(awsconfig *aws.Config, cfg *flags.API) Customers {
	session := dynastore.New(awsconfig)

	tbl := session.Table(cfg.CustomersTable)

	return &DDBCustomers{
		customerPart:      tbl.Partition(customersPartitionName),
		customerNamesPart: tbl.Partition(customerNamesPartitionName),
	}
}
