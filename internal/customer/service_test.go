package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
	mocks "github.com/danilosano/web-golang-api/pkg/tests/customers"
	"github.com/stretchr/testify/assert"
)

func createService(t *testing.T) (Service, *mocks.CustomersRepositoryMock, context.Context) {
	t.Helper()
	repoMock := new(mocks.CustomersRepositoryMock)
	service := NewService(repoMock)
	ctx := context.Background()
	return service, repoMock, ctx
}

var (
	CustomerNumber int = 2

	mockedCustomerList = []dto.ResultCustomerRequest{
		{
			ID:             1,
			CustomerNumber: &CustomerNumber,
			FirstName:      "Danilo",
			LastName:       "Sano",
			CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
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

	mockedResultCustomer = dto.ResultCustomerRequest{
		ID:             1,
		CustomerNumber: &CustomerNumber,
		FirstName:      mockedCustomer.FirstName,
		LastName:       mockedCustomer.LastName,
	}
)

func TestGetAll(t *testing.T) {
	t.Run(`If the list has "n" elements, it will return an amount of the total elements.`, func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("GetAllWithContext", ctx).Return(mockedCustomerList, nil)

		result, err := service.GetAll(ctx)
		assert.Nil(t, err)
		assert.True(t, len(result) > 0)
	})

	t.Run("When the backend returns an unexpected error, return the error.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("GetAllWithContext", ctx).Return([]domain.Customer{}, errors.New("generic error"))

		_, err := service.GetAll(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("If it contains the required fields, it will be created.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("SaveWithContext", ctx, domain.Customer{
			ID:             0,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			CreatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(1, nil)
		repoMock.On("GetWithContext", ctx, 1).Return(mockedResultCustomer, nil)

		result, err := service.Save(ctx, input)
		assert.Nil(t, err)
		mockedResultCustomer.CreatedAt = result.CreatedAt
		assert.Equal(t, mockedResultCustomer, result)
	})

	t.Run("If the customer_number already exists it cannot be created.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(true)

		_, err := service.Save(ctx, input)
		assert.NotNil(t, err)
		assert.Equal(t, ErrorCustomerNumberAlreadyExist, err)
	})

	t.Run("If an unexpected backend error occurs in the save function, return an error.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("SaveWithContext", ctx, domain.Customer{
			ID:             0,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			CreatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(1, errors.New("generic error"))

		_, err := service.Save(ctx, input)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})

	t.Run("If an unexpected error occurs on the backend when trying to retrieve the object, return an error.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("SaveWithContext", ctx, domain.Customer{
			ID:             0,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			CreatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(1, nil)
		repoMock.On("GetWithContext", ctx, 1).Return(domain.Customer{}, errors.New("generic error"))

		_, err := service.Save(ctx, input)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})
}

func TestGet(t *testing.T) {
	t.Run("If the element searched for by id does not exist, return null.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(false)

		_, err := service.Get(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, ErrorCustomerNotFound, err)
	})

	t.Run("If the element searched for by id exists, it will return the requested element information.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(true)
		repoMock.On("GetWithContext", ctx, 1).Return(mockedResultCustomer, nil)

		result, err := service.Get(ctx, 1)
		assert.Nil(t, err)
		assert.Equal(t, mockedResultCustomer, result)
	})

	t.Run("When the backend returns an unexpected error, return the error", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(true)
		repoMock.On("GetWithContext", ctx, 1).Return([]domain.Customer{}, errors.New("generic error"))

		_, err := service.Get(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("When the data update is successful, the customer with the updated information will be returned.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByIDWithContext", ctx, mockedCustomer.ID).Return(true)
		repoMock.On("ExistsByCustomerNumberAndIDWithContext", ctx, mockedCustomer.ID, *input.CustomerNumber).Return(true)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("UpdateWithContext", ctx, domain.Customer{
			ID:             mockedCustomer.ID,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			UpdatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(nil)
		repoMock.On("GetWithContext", ctx, mockedCustomer.ID).Return(mockedResultCustomer, nil)

		result, err := service.Update(ctx, inputUpdate, mockedCustomer.ID)
		assert.Nil(t, err)
		assert.Equal(t, mockedResultCustomer, result)
	})

	t.Run("If the customer to be updated does not exist, null will be returned.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByIDWithContext", ctx, mockedCustomer.ID).Return(false)

		_, err := service.Update(ctx, inputUpdate, mockedCustomer.ID)
		assert.NotNil(t, err)
		assert.Equal(t, ErrorCustomerNotFound, err)
	})

	t.Run("If the backend returns an unexpected error when trying to update, return the error.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByIDWithContext", ctx, mockedCustomer.ID).Return(true)
		repoMock.On("ExistsByCustomerNumberAndIDWithContext", ctx, mockedCustomer.ID, *input.CustomerNumber).Return(true)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("UpdateWithContext", ctx, domain.Customer{
			ID:             mockedCustomer.ID,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			UpdatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(errors.New("generic error"))

		_, err := service.Update(ctx, inputUpdate, mockedCustomer.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})

	t.Run("If the backend returns an unexpected error when returning the object, return the error.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByIDWithContext", ctx, mockedCustomer.ID).Return(true)
		repoMock.On("ExistsByCustomerNumberAndIDWithContext", ctx, mockedCustomer.ID, *input.CustomerNumber).Return(true)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(false)
		repoMock.On("UpdateWithContext", ctx, domain.Customer{
			ID:             mockedCustomer.ID,
			CustomerNumber: *input.CustomerNumber,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			UpdatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Now().Location())}).Return(nil)
		repoMock.On("GetWithContext", ctx, mockedCustomer.ID).Return(domain.Customer{}, errors.New("generic error"))

		_, err := service.Update(ctx, inputUpdate, mockedCustomer.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})

	t.Run("If the customer to be updated exists, but the customer number to be updated already exists.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)
		repoMock.On("ExistsByIDWithContext", ctx, mockedCustomer.ID).Return(true)
		repoMock.On("ExistsByCustomerNumberAndIDWithContext", ctx, mockedCustomer.ID, *input.CustomerNumber).Return(false)
		repoMock.On("ExistsByCustomerNumberWithContext", ctx, *input.CustomerNumber).Return(true)

		_, err := service.Update(ctx, inputUpdate, mockedCustomer.ID)
		assert.NotNil(t, err)
		assert.Equal(t, ErrorCustomerNumberAlreadyExist, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("When the customer does not exist, null will be returned.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(false)

		err := service.Delete(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, ErrorCustomerNotFound, err)
	})

	t.Run("If the deletion is successful, the item will not appear in the list.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(true)
		repoMock.On("DeleteWithContext", ctx, 1).Return(nil)

		err := service.Delete(ctx, 1)
		assert.Nil(t, err)
	})

	t.Run("If an error occurs while deleting, the error will be returned.", func(t *testing.T) {
		service, repoMock, ctx := createService(t)

		repoMock.On("ExistsByIDWithContext", ctx, 1).Return(true)
		repoMock.On("DeleteWithContext", ctx, 1).Return(errors.New("generic error"))

		err := service.Delete(ctx, 1)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("generic error"), err)
	})
}
