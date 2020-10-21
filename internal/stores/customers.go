package stores

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/dynastore"
)

const customerPrefix = "CUST"

type DDBCustomers struct {
	customerPart dynastore.Partition
}

func (dc *DDBCustomers) GetCustomer(ctx context.Context, id string, into proto.Message) (int64, error) {

	log.Ctx(ctx).Info().Str("id", id).Msg("get customer")

	key := fmt.Sprintf("%s_%s", customerPrefix, id)

	kv, err := dc.customerPart.Get(key)
	if err != nil {
		return 0, err
	}

	if err := proto.Unmarshal(kv.BytesValue(), into); err != nil {
		return 0, err
	}

	return kv.Version, nil
}

func (dc *DDBCustomers) CreateCustomer(ctx context.Context, id string, obj proto.Message) (int64, error) {

	log.Ctx(ctx).Info().Str("id", id).Msg("create customer")

	data, err := proto.Marshal(obj)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s_%s", customerPrefix, id)

	_, kv, err := dc.customerPart.AtomicPut(key, dynastore.WriteWithBytes(data), dynastore.WriteWithNoExpires())
	if err != nil {
		return 0, err
	}

	return kv.Version, nil
}
