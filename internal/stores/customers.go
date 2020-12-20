package stores

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/dynastore"
	"google.golang.org/protobuf/proto"
)

const (
	customersPartitionName     = "customers"
	customerNamesPartitionName = "customer_names"
	customerRecordPrefix       = "CUST_"
)

var (
	// ErrCustomerNameConflict when creating a new customer a name conflict was found
	ErrCustomerNameConflict = errors.New("customer name already exists")
	// ErrCustomerNotFound when retrieving a customer the id wasn't found
	ErrCustomerNotFound = errors.New("customer not found")
)

// DDBCustomers DynamoDB based customer store
type DDBCustomers struct {
	customerPart      dynastore.Partition
	customerNamesPart dynastore.Partition
}

// GetCustomer retrieve the customer and marshall it into the supplied message
func (dc *DDBCustomers) GetCustomer(ctx context.Context, id string, into proto.Message) (int64, error) {
	log.Ctx(ctx).Info().Str("id", id).Msg("get customer")

	key := fmt.Sprintf("%s%s", customerRecordPrefix, id)

	kv, err := dc.customerPart.Get(key)
	if err != nil {
		if err == dynastore.ErrKeyNotFound {
			return 0, ErrCustomerNotFound
		}
		return 0, err
	}

	if err := proto.Unmarshal(kv.BytesValue(), into); err != nil {
		return 0, err
	}

	return kv.Version, nil
}

// ListCustomers list customers and retrieve the raw records for marshaling in the consuming method, this is done to avoid
// tying this method to a given storage model
func (dc *DDBCustomers) ListCustomers(ctx context.Context, nextToken string, limit int) (string, []Record, error) {
	log.Ctx(ctx).Info().Msg("list customers")

	readOpts := []dynastore.ReadOption{}

	if nextToken != "" {
		readOpts = append(readOpts, dynastore.ReadWithStartKey(nextToken))
	}

	if limit != 0 {
		readOpts = append(readOpts, dynastore.ReadWithLimit(int64(limit)))
	}

	kvpage, err := dc.customerPart.ListPage(customerRecordPrefix, readOpts...)
	if err != nil {
		return "", nil, err
	}

	records := make([]Record, len(kvpage.Keys))

	for i, kv := range kvpage.Keys {
		records[i] = Record{
			Data:    kv.BytesValue(),
			Version: kv.Version,
			ID:      strings.TrimPrefix(kv.Key, customerRecordPrefix),
		}
	}

	return kvpage.LastKey, records, nil
}

// CreateCustomer create a customer using the supplied message and return the version
func (dc *DDBCustomers) CreateCustomer(ctx context.Context, id, name string, obj proto.Message) (int64, error) {
	log.Ctx(ctx).Info().Str("id", id).Str("name", name).Msg("create customer")

	exists, err := dc.customerNamesPart.Exists(name)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrCustomerNameConflict
	}

	// AtomicPut will return an error if the key already exists
	_, kv, err := dc.customerNamesPart.AtomicPut(name, dynastore.WriteWithNoExpires(), dynastore.WriteWithString(id))
	if err != nil {
		return 0, err
	}

	log.Ctx(ctx).Info().Str("name", name).Int64("v", kv.Version).Msg("name created")

	data, err := proto.Marshal(obj)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s%s", customerRecordPrefix, id)

	_, kv, err = dc.customerPart.AtomicPut(key, dynastore.WriteWithBytes(data), dynastore.WriteWithNoExpires())
	if err != nil {
		return 0, err
	}

	return kv.Version, nil
}
