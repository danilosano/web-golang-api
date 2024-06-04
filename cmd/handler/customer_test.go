package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/danilosano/web-golang-api/internal/customer"
	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
	mocks "github.com/danilosano/web-golang-api/pkg/tests/customers"
	"github.com/danilosano/web-golang-api/pkg/testutil"
	"github.com/danilosano/web-golang-api/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	pathCustomer = "/api/v1/customers/"
)

func InitServerWithCustomersRoute(t *testing.T) (*gin.Engine, *mocks.CustomersServiceMock, context.Context) {
	t.Helper()
	server := testutil.CreateServer()
	mockService := new(mocks.CustomersServiceMock)
	handler := NewCustomerHandler(mockService)
	server.GET(pathCustomer, handler.GetAll)
	server.GET(pathCustomer+":id", handler.Get)
	server.POST(pathCustomer, handler.Store)
	server.PUT(pathCustomer+":id", handler.Update)
	server.DELETE(pathCustomer+":id", handler.Delete)
	ctx := context.Background()
	return server, mockService, ctx
}

var (
	CustomerNumber int = 2
	mockedCustomer     = domain.Customer{
		ID:             1,
		CustomerNumber: 2,
		FirstName:      "Danilo",
		LastName:       "Sano",
		CreatedAt:      time.Date(2021, 10, 10, 0, 0, 0, 0, time.UTC),
	}

	mockedResultCustomer = dto.ResultCustomerRequest{
		ID:             1,
		CustomerNumber: &CustomerNumber,
		FirstName:      "Danilo",
		LastName:       "Sano",
		CreatedAt:      time.Date(2021, 10, 10, 0, 0, 0, 0, time.UTC),
	}

	input = dto.CreateCustomerRequest{
		CustomerNumber: &CustomerNumber,
		FirstName:      "Danilo",
		LastName:       "Sano",
	}

	inputUpdate = dto.UpdateCustomerRequest{
		CustomerNumber: &CustomerNumber,
		FirstName:      "Danilo",
		LastName:       "Sano",
	}

	jsonInput = `{
		"customer_number":2,
		"first_name":"Danilo",
		"last_name":"Sano"
	}`
)

func TestCreate(t *testing.T) {
	t.Run("When data entry is successful, a 201 code will be returned along with the inserted object.", func(t *testing.T) {
		var data dto.ResultCustomerRequest
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Save", ctx, input).Return(mockedResultCustomer, nil)

		request, response := testutil.MakeRequest(http.MethodPost, pathCustomer, jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &data)
		assert.Nil(t, err)
		assert.Equal(t, mockedResultCustomer, data)
	})

	t.Run("If the JSON object does not contain the required fields, a 422 code will be returned.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Save", ctx, input).Return(domain.Customer{}, nil)

		request, response := testutil.MakeRequest(http.MethodPost, pathCustomer, `{
			"customer_number":1,
			"first_name":,
			"last_name":"Sano"
		}`)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("If the customer_number already exists, it will return a 409 Conflict error.", func(t *testing.T) {
		var resp web.ErrorResponse
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Save", ctx, input).Return(domain.Customer{}, customer.ErrorCustomerNumberAlreadyExist)

		request, response := testutil.MakeRequest(http.MethodPost, pathCustomer, jsonInput)
		server.ServeHTTP(response, request)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusConflict, response.Code)
		assert.Equal(t, customer.ErrorCustomerNumberAlreadyExist.Error(), resp.Message)
	})

	t.Run("When data entry is successful, but an internal server error occurs when creating.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Save", ctx, input).Return(domain.Customer{}, errors.New("generic error"))

		request, response := testutil.MakeRequest(http.MethodPost, pathCustomer, jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("When the request is successful, the backend returns a list of all existing customers.", func(t *testing.T) {
		var result web.Responses
		var data []dto.ResultCustomerRequest

		mockedCustomersList := []dto.ResultCustomerRequest{mockedResultCustomer}
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("GetAll", ctx).Return(mockedCustomersList, nil)

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer, "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &result)
		assert.Nil(t, err)
		jsonData, err := json.Marshal(result.Data)
		assert.Nil(t, err)
		err = json.Unmarshal(jsonData, &data)
		assert.Nil(t, err)
		assert.Equal(t, mockedCustomersList, data)
	})

	t.Run("When an unexpected error occurs in the backend, it will return a 500 error.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("GetAll", ctx).Return([]domain.Customer{}, errors.New("generic error"))

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer, "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestGet(t *testing.T) {
	t.Run("When the customer does not exist, a 404 code will be returned.", func(t *testing.T) {
		var result web.ErrorResponse

		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Get", ctx, 1).Return(dto.ResultCustomerRequest{}, customer.ErrorCustomerNotFound)

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer+"1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, customer.ErrorCustomerNotFound.Error(), result.Message)
	})

	t.Run("When the request is successful, the backend will return the requested customer information.", func(t *testing.T) {
		var result web.Responses
		var data dto.ResultCustomerRequest

		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Get", ctx, 1).Return(mockedResultCustomer, nil)

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer+"1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &result)
		assert.Nil(t, err)
		jsonData, err := json.Marshal(result.Data)
		assert.Nil(t, err)
		err = json.Unmarshal(jsonData, &data)
		assert.Nil(t, err)
		assert.Equal(t, mockedResultCustomer, data)
	})

	t.Run("When the backend returns an unexpected error, return code 500.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Get", ctx, 1).Return(domain.Customer{}, errors.New("generic error"))

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer+"1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("When the parameter id is wrong, it will return a Bad Request status.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, _, _ := InitServerWithCustomersRoute(t)

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer+"-1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid input ID", resp.Message)
	})

	t.Run("When the parameter id is zero, a 400 Bad Request code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, _, _ := InitServerWithCustomersRoute(t)

		request, response := testutil.MakeRequest(http.MethodGet, pathCustomer+"0", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid id provided: id must be a positive non-zero number", resp.Message)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("When the data update is successful, the customer with the updated information will be returned along with a 200 code.", func(t *testing.T) {
		var data dto.ResultCustomerRequest

		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Update", ctx, inputUpdate, 1).Return(mockedResultCustomer, nil)

		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"1", jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &data)
		assert.Nil(t, err)
		assert.Equal(t, mockedResultCustomer, data)
	})

	t.Run("When the backend returns an unexpected error, return code 500.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Update", ctx, input, 1).Return(domain.Customer{}, errors.New("generic error"))

		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"1", jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("If the customer to be updated does not exist, a 404 code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Update", ctx, inputUpdate, 9999).Return(domain.Customer{}, customer.ErrorCustomerNotFound)

		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"9999", jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, customer.ErrorCustomerNotFound.Error(), resp.Message)
	})

	t.Run("When the parameter id is wrong, it will return a Bad Request status.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, _, _ := InitServerWithCustomersRoute(t)

		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"-1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid input ID", resp.Message)
	})

	t.Run("When a field is missing in the request body, a 422 code will be returned.", func(t *testing.T) {
		server, _, _ := InitServerWithCustomersRoute(t)
		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"1", `{}`)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("When there is an invalid value in the request body, a 400 code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse
		server, _, _ := InitServerWithCustomersRoute(t)
		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"1", `{
			"customer_number":-1,
			"first_name":"aaa",
			"last_name":"Sano"
		}`)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid input: customer number must be greater than 0", resp.Message)
	})

	t.Run("If the customer number to be updated already exists, a 409 code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Update", ctx, inputUpdate, 2).Return(domain.Customer{}, customer.ErrorCustomerNumberAlreadyExist)

		request, response := testutil.MakeRequest(http.MethodPut, pathCustomer+"2", jsonInput)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, customer.ErrorCustomerNumberAlreadyExist.Error(), resp.Message)
	})
}

func TestDelete(t *testing.T) {
	t.Run("When the customer does not exist, a 404 code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Delete", ctx, 9999).Return(customer.ErrorCustomerNotFound)

		request, response := testutil.MakeRequest(http.MethodDelete, pathCustomer+"9999", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, customer.ErrorCustomerNotFound.Error(), resp.Message)
	})

	t.Run("When the deletion is successful, a 204 code will be returned.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Delete", ctx, 1).Return(nil)

		request, response := testutil.MakeRequest(http.MethodDelete, pathCustomer+"1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("When the parameter id is wrong, a Bad Request code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, _, _ := InitServerWithCustomersRoute(t)

		request, response := testutil.MakeRequest(http.MethodDelete, pathCustomer+"-1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid ID provided", resp.Message)
	})

	t.Run("When the parameter id is zero, a 400 Bad Request code will be returned.", func(t *testing.T) {
		var resp web.ErrorResponse

		server, _, _ := InitServerWithCustomersRoute(t)

		request, response := testutil.MakeRequest(http.MethodDelete, pathCustomer+"0", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		err := json.Unmarshal(response.Body.Bytes(), &resp)
		assert.Nil(t, err)
		assert.Equal(t, "invalid id provided: id must be a positive non-zero number", resp.Message)
	})

	t.Run("When the backend returns an unexpected error, return code 500.", func(t *testing.T) {
		server, service, ctx := InitServerWithCustomersRoute(t)
		service.On("Delete", ctx, 1).Return(domain.Customer{}, errors.New("generic error"))

		request, response := testutil.MakeRequest(http.MethodDelete, pathCustomer+"1", "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
