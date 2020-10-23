package stores

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/dynastore"
)

const (
	customersPartitionName     = "customers"
	customerNamesPartitionName = "customer_names"
	customerRecordPrefix       = "CUST_"
)

var (
	ErrCustomerNameConfict = errors.New("customer name already exists")
	ErrCustomerNotFound    = errors.New("customer not found")
)

type DDBCustomers struct {
	customerPart      dynastore.Partition
	customerNamesPart dynastore.Partition
}

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

func (dc *DDBCustomers) CreateCustomer(ctx context.Context, id, name string, obj proto.Message) (int64, error) {

	log.Ctx(ctx).Info().Str("id", id).Str("name", name).Msg("create customer")

	exists, err := dc.customerNamesPart.Exists(name)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrCustomerNameConfict
	}

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
