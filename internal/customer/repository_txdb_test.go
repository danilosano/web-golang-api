package customer_test

import (
	"context"
	"testing"
	"time"

	"github.com/danilosano/web-golang-api/internal/customer"
	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
	"github.com/danilosano/web-golang-api/pkg/testutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	mockedCustomer = domain.Customer{
		CustomerNumber: 999,
		FirstName:      "Danilo",
		LastName:       "Sano",
		CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockedCustomerUpdated = domain.Customer{
		CustomerNumber: 999,
		FirstName:      "Danilo",
		LastName:       "Sano",
		CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
)

func TestSuite_CustomerRepository(t *testing.T) {
	db, err := testutil.InitTxdbDatabase(t)
	assert.NoError(t, err)
	repository := customer.NewRepository(db)

	testStoreWithContext(t, repository)
	testUpdateWithContext(t, repository)
	testDeleteWithContext(t, repository)
	testExistsByCustomerNumberWithContext(t, repository)
	testExistsByCustomerNumberAndIDWithContext(t, repository)
	testGetByCustomerNumberWithContext(t, repository)
	testGetAllWithContext(t, repository)

	db.Close()
}

func testStoreWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := repository.SaveWithContext(ctx, mockedCustomer)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	result := repository.ExistsByIDWithContext(ctx, id)
	assert.True(t, result)
}

func testUpdateWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := repository.SaveWithContext(ctx, mockedCustomer)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	exists := repository.ExistsByIDWithContext(ctx, id)
	assert.True(t, exists)

	mockedCustomerUpdated.ID = id

	err = repository.UpdateWithContext(ctx, mockedCustomerUpdated)
	assert.NoError(t, err)

	result, err := repository.GetWithContext(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, dto.ResultCustomerRequest{
		ID:             id,
		CustomerNumber: &mockedCustomerUpdated.CustomerNumber,
		FirstName:      mockedCustomerUpdated.FirstName,
		LastName:       mockedCustomerUpdated.LastName,
		CreatedAt:      mockedCustomer.CreatedAt,
		UpdatedAt:      &mockedCustomerUpdated.UpdatedAt,
	}, result)
}

func testDeleteWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := repository.SaveWithContext(ctx, mockedCustomer)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	exists := repository.ExistsByIDWithContext(ctx, id)
	assert.True(t, exists)

	err = repository.DeleteWithContext(ctx, id)
	assert.NoError(t, err)

	exists = repository.ExistsByIDWithContext(ctx, id)
	assert.False(t, exists)
}

func testExistsByCustomerNumberWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := repository.SaveWithContext(ctx, mockedCustomer)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	exists := repository.ExistsByCustomerNumberWithContext(ctx, mockedCustomer.CustomerNumber)
	assert.True(t, exists)
}

func testExistsByCustomerNumberAndIDWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := repository.SaveWithContext(ctx, mockedCustomer)
	assert.NoError(t, err)
	assert.True(t, id > 0)

	exists := repository.ExistsByCustomerNumberAndIDWithContext(ctx, id, mockedCustomer.CustomerNumber)
	assert.True(t, exists)
}

func testGetByCustomerNumberWithContext(t *testing.T, repository customer.Repository) {
	t.Run("If the customer number exists, return the customer", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		id, err := repository.SaveWithContext(ctx, mockedCustomer)
		assert.NoError(t, err)
		assert.True(t, id > 0)

		mockedCustomer.ID = id

		result, err := repository.GetByCustomerNumberWithContext(ctx, mockedCustomer.CustomerNumber)
		assert.NoError(t, err)
		assert.Equal(t, mockedCustomer, domain.Customer{
			ID:             id,
			CustomerNumber: *result.CustomerNumber,
			FirstName:      result.FirstName,
			LastName:       result.LastName,
			CreatedAt:      result.CreatedAt,
		})
	})

	t.Run("If the customer number does not exists, return a not found error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := repository.GetByCustomerNumberWithContext(ctx, 999999)
		assert.Error(t, err)
		assert.Equal(t, customer.ErrorCustomerNotFound, err)
	})
}

func testGetAllWithContext(t *testing.T, repository customer.Repository) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	report, err := repository.GetAllWithContext(ctx)
	assert.NoError(t, err)
	assert.True(t, len(report) > 0)
}
