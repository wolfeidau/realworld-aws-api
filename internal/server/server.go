package server

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/stores"
	storagepb "github.com/wolfeidau/realworld-aws-api/proto/customers/storage/v1beta1"
)

// Customers customers REST API
type Customers struct {
	cfg           *flags.API
	customerStore stores.Customers
}

// NewCustomers construct a new customers with the supplied configuration
func NewCustomers(cfg *flags.API, customerStore stores.Customers) *Customers {
	return &Customers{
		cfg:           cfg,
		customerStore: customerStore,
	}
}

// Customers Get a list of customers.
// (GET /customers)
func (cs *Customers) Customers(c echo.Context, params customersapi.CustomersParams) error {
	ctx := c.Request().Context()

	nextToken, records, err := cs.customerStore.ListCustomers(ctx, aws.StringValue(params.NextToken), aws.IntValue(params.MaxItems))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	customersPage, err := toAPICustomersPage(records, nextToken)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert customer page")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &customersPage)
}

// NewCustomer Create a customer.
// (POST /customers)
func (cs *Customers) NewCustomer(c echo.Context) error {
	ctx := c.Request().Context()

	newCust := new(customersapi.NewCustomer)

	err := c.Bind(newCust)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to parse new customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	id := mustNewID()
	storedCust := toStorageCustomer(newCust)

	v, err := cs.customerStore.CreateCustomer(ctx, id, newCust.Name, storedCust)
	if err != nil {
		if err == stores.ErrCustomerNameConfict {
			return c.NoContent(http.StatusConflict)
		}
		log.Ctx(ctx).Error().Err(err).Msg("failed to store new customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.Ctx(ctx).Info().Str("id", id).Int64("v", v).Msg("stored new customer")

	cust, err := fromStorageCustomer(id, storedCust)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to marshal new customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, cust)
}

// GetCustomer (GET /customers/{id})
func (cs *Customers) GetCustomer(c echo.Context, id string) error {
	ctx := c.Request().Context()

	cust := &storagepb.Customer{}

	_, err := cs.customerStore.GetCustomer(c.Request().Context(), id, cust)
	if err != nil {
		if err == stores.ErrCustomerNotFound {
			return c.NoContent(http.StatusNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Msg("failed to get customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	apiCust, err := fromStorageCustomer(id, cust)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to marshal customer")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, apiCust)
}

// UpdateCustomer Update a customer.
// (PUT /customers/{id})
func (cs *Customers) UpdateCustomer(c echo.Context, id string) error {
	return c.NoContent(http.StatusNotImplemented)
}
