package stores

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/logger"
	storagepb "github.com/wolfeidau/realworld-aws-api/proto/customers/storage/v1beta1"
)

func TestDDBCustomers(t *testing.T) {
	assert := require.New(t)

	err := ensureVersionTable(dbSvc, "customers-test-table")
	assert.NoError(err)

	awscfg := mustConfig(endpoint)

	ctx := logger.NewLoggerWithContext(context.TODO())

	customers := NewCustomers(awscfg, &flags.API{CustomersTable: "customers-test-table"})

	t.Run("create", func(t *testing.T) {
		id, name, labels := "abc123", "new customer", []string{"test"}
		cust := newStoreCustomer(name, labels)

		v, err := customers.CreateCustomer(ctx, id, name, cust)
		assert.NoError(err)
		assert.Equal(int64(1), v)

		// should conflict
		v, err = customers.CreateCustomer(ctx, id, name, cust)
		assert.Equal(ErrCustomerNameConfict, err)
		assert.Equal(int64(0), v)
	})

	t.Run("get", func(t *testing.T) {

		id, name, labels := "cde456", "new get customer", []string{"test"}
		newCust := newStoreCustomer(name, labels)

		v, err := customers.CreateCustomer(ctx, id, name, newCust)
		assert.NoError(err)
		assert.Equal(int64(1), v)

		cust := new(storagepb.Customer)

		v, err = customers.GetCustomer(ctx, id, cust)
		assert.NoError(err)
		assert.Equal(int64(1), v)
		assert.Equal(name, cust.Name)
	})

}

func newStoreCustomer(name string, labels []string) *storagepb.Customer {
	return &storagepb.Customer{
		Name:    name,
		Labels:  labels,
		Created: ptypes.TimestampNow(),
		Updated: ptypes.TimestampNow(),
	}
}
