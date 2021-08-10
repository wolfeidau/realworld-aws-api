package server

import (
	"crypto/rand"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/oklog/ulid/v2"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/stores"
	storagepb "github.com/wolfeidau/realworld-aws-api/proto/customers/storage/v1beta1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toStorageCustomer(newCust *customersapi.NewCustomer) *storagepb.Customer {
	cust := &storagepb.Customer{
		Name:    newCust.Name,
		Labels:  newCust.Labels,
		Created: timestamppb.Now(),
		Updated: timestamppb.Now(),
	}

	if newCust.Description != nil {
		cust.Description = &wrappers.StringValue{Value: *newCust.Description}
	}

	return cust
}

func fromStorageCustomer(id string, cust *storagepb.Customer) *customersapi.Customer {
	resCust := &customersapi.Customer{
		Id:        id,
		Name:      cust.Name,
		Labels:    cust.Labels,
		CreatedAt: cust.Created.AsTime(),
		UpdatedAt: cust.Updated.AsTime(),
	}

	if cust.Description != nil {
		resCust.Description = &cust.Description.Value
	}

	return resCust
}

func toAPICustomersPage(records []stores.Record, nextToken string) (*customersapi.CustomersPage, error) {
	customersPage := new(customersapi.CustomersPage)
	customersPage.Customers = make([]customersapi.Customer, len(records))

	if nextToken != "" {
		customersPage.NextToken = &nextToken
	}

	for i, record := range records {
		storedCust := &storagepb.Customer{}

		err := proto.Unmarshal(record.Data, storedCust)
		if err != nil {
			return nil, err
		}

		cust := fromStorageCustomer(record.ID, storedCust)
		if err != nil {
			return nil, err
		}

		customersPage.Customers[i] = *cust
	}

	return customersPage, nil
}

func mustNewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}
