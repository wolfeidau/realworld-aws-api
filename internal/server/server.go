package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
)

// Customers customers REST API
type Customers struct {
}

// NewCustomers construct a new customers with the supplied configuration
func NewCustomers(cfg *flags.API) *Customers {
	return &Customers{}
}

// Customers Get a list of customers.
// (GET /customers)
func (cs *Customers) Customers(c echo.Context, params customersapi.CustomersParams) error {
	return c.NoContent(http.StatusNotImplemented)
}

// NewCustomer Create a customer.
// (POST /customers)
func (cs *Customers) NewCustomer(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

// GetCustomer (GET /customers/{id})
func (cs *Customers) GetCustomer(c echo.Context, id string) error {
	return c.NoContent(http.StatusNotImplemented)
}

// UpdateCustomer Update a customer.
// (PUT /customers/{id})
func (cs *Customers) UpdateCustomer(c echo.Context, id string) error {
	return c.NoContent(http.StatusNotImplemented)
}
