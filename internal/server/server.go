package server

import (
	"net/http"
	"time"

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
	return c.NoContent(http.StatusNotImplemented)
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
	storedCust := createStorageCustomer(newCust)

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

	start := time.Now()

	cust := &storagepb.Customer{}

	v, err := cs.customerStore.GetCustomer(c.Request().Context(), id, cust)
	if err != nil {
		return c.JSON(500, map[string]string{"msg": "failed to get customer"})
	}

	apiCust, err := fromStorageCustomer(id, cust)
	if err != nil {
		return c.JSON(500, map[string]string{"msg": "failed to get customer"})
	}

	log.Ctx(c.Request().Context()).Info().Dur("duration", time.Since(start)).Int64("v", v).Msg("get customer")

	return c.JSON(http.StatusOK, apiCust)
}

// UpdateCustomer Update a customer.
// (PUT /customers/{id})
func (cs *Customers) UpdateCustomer(c echo.Context, id string) error {
	return c.NoContent(http.StatusNotImplemented)
}
