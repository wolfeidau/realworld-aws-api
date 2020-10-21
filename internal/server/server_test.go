package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"

	"github.com/golang/mock/gomock"
	"github.com/wolfeidau/realworld-aws-api/mocks"

	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
)

func TestCustomers_NewCustomer(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerStore := mocks.NewMockCustomers(ctrl)

	customerStore.EXPECT().CreateCustomer(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)

	cs := Customers{cfg: &flags.API{}, customerStore: customerStore}

	e := echo.New()

	zlog := zerolog.New(os.Stderr).With().
		Stack().Caller().Logger()

	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(`{"labels":["test"],"name":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(zlog.WithContext(context.TODO()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := cs.NewCustomer(c)
	assert.NoError(err)
	assert.Equal(http.StatusCreated, rec.Code)

	cust := new(customersapi.Customer)
	err = json.Unmarshal(rec.Body.Bytes(), cust)
	assert.NoError(err)

	assert.Equal("test", cust.Name)
	assert.Equal([]string{"test"}, cust.Labels)

}
