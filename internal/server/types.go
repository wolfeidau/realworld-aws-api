package server

import (
	"crypto/rand"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/oklog/ulid/v2"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	storagepb "github.com/wolfeidau/realworld-aws-api/proto/customers/storage/v1beta1"
)

func createStorageCustomer(newCust *customersapi.NewCustomer) *storagepb.Customer {
	cust := &storagepb.Customer{
		Name:    newCust.Name,
		Labels:  newCust.Labels,
		Created: ptypes.TimestampNow(),
		Updated: ptypes.TimestampNow(),
	}

	if newCust.Description != nil {
		cust.Description = &wrappers.StringValue{Value: *newCust.Description}
	}

	return cust
}

func fromStorageCustomer(id string, cust *storagepb.Customer) (*customersapi.Customer, error) {

	created, err := ptypes.Timestamp(cust.Created)
	if err != nil {
		return nil, err
	}
	updated, err := ptypes.Timestamp(cust.Updated)
	if err != nil {
		return nil, err
	}

	resCust := &customersapi.Customer{
		Id:        id,
		Name:      cust.Name,
		Labels:    cust.Labels,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	if cust.Description != nil {
		resCust.Description = &cust.Description.Value
	}

	return resCust, nil
}

func mustNewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}
