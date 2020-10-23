package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/realworld-aws-api/internal/customersapi"
	"github.com/wolfeidau/realworld-aws-api/internal/flags"
	"github.com/wolfeidau/realworld-aws-api/internal/logger"
	"github.com/wolfeidau/realworld-aws-api/mocks"
)

func TestCustomers_NewCustomer(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerStore := mocks.NewMockCustomers(ctrl)

	customerStore.EXPECT().CreateCustomer(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)

	cs := Customers{cfg: &flags.API{}, customerStore: customerStore}

	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(`{"labels":["test"],"name":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(logger.NewLoggerWithContext(context.TODO()))
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

func TestCustomers_GetCustomer(t *testing.T) {
	assert := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	customerStore := mocks.NewMockCustomers(ctrl)

	callbackFunc := func(ctx context.Context, id string, into proto.Message) (int64, error) {
		data, err := base64.StdEncoding.DecodeString("CgR0ZXN0EgYKBHRlc3QaBHRlc3QiCwiah8f8BRDG++FzKgsImofH/AUQ/4Dicw==")
		assert.NoError(err)

		if err := proto.Unmarshal(data, into); err != nil {
			return 0, err
		}

		return 1, nil
	}

	customerStore.EXPECT().GetCustomer(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(callbackFunc)

	cs := Customers{cfg: &flags.API{}, customerStore: customerStore}

	e := echo.New()

	id := "abc123"

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/customers/%s", id), nil)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(logger.NewLoggerWithContext(context.TODO()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := cs.GetCustomer(c, id)
	assert.NoError(err)
	assert.Equal(http.StatusOK, rec.Code)

	cust := new(customersapi.Customer)
	err = json.Unmarshal(rec.Body.Bytes(), cust)
	assert.NoError(err)
	assert.Equal("test", cust.Name)
	assert.Equal([]string{"test"}, cust.Labels)

}
