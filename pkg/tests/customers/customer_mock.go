package mocks

import (
	"context"

	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
	"github.com/stretchr/testify/mock"
)

type CustomersServiceMock struct {
	mock.Mock
}

func (p *CustomersServiceMock) GetAll(ctx context.Context) ([]dto.ResultCustomerRequest, error) {
	args := p.Called(ctx)

	arg0, ok := args.Get(0).([]dto.ResultCustomerRequest)
	if !ok {
		return []dto.ResultCustomerRequest{}, args.Error(1)

	}
	return arg0, args.Error(1)
}

func (p *CustomersServiceMock) Save(ctx context.Context, s dto.CreateCustomerRequest) (dto.ResultCustomerRequest, error) {
	args := p.Called(ctx, s)

	arg0, ok := args.Get(0).(dto.ResultCustomerRequest)
	if !ok {
		return dto.ResultCustomerRequest{}, args.Error(1)

	}

	return arg0, args.Error(1)
}

func (p *CustomersServiceMock) Delete(ctx context.Context, id int) error {
	args := p.Called(ctx, id)
	return args.Error(0)
}

func (p *CustomersServiceMock) Update(ctx context.Context, s dto.UpdateCustomerRequest, id int) (dto.ResultCustomerRequest, error) {
	args := p.Called(ctx, s, id)

	arg0, ok := args.Get(0).(dto.ResultCustomerRequest)
	if !ok {
		return dto.ResultCustomerRequest{}, args.Error(1)

	}

	return arg0, args.Error(1)
}

func (p *CustomersServiceMock) Get(ctx context.Context, id int) (dto.ResultCustomerRequest, error) {
	args := p.Called(ctx, id)

	arg0, ok := args.Get(0).(dto.ResultCustomerRequest)
	if !ok {
		return dto.ResultCustomerRequest{}, args.Error(1)

	}

	return arg0, args.Error(1)
}

type CustomersRepositoryMock struct {
	mock.Mock
}

func (s *CustomersRepositoryMock) GetAllWithContext(ctx context.Context) ([]dto.ResultCustomerRequest, error) {
	args := s.Called(ctx)

	arg0, ok := args.Get(0).([]dto.ResultCustomerRequest)
	if !ok {
		return []dto.ResultCustomerRequest{}, args.Error(1)

	}
	return arg0, args.Error(1)
}

func (s *CustomersRepositoryMock) GetWithContext(ctx context.Context, id int) (dto.ResultCustomerRequest, error) {
	args := s.Called(ctx, id)

	arg0, ok := args.Get(0).(dto.ResultCustomerRequest)
	if !ok {
		return dto.ResultCustomerRequest{}, args.Error(1)
	}

	return arg0, args.Error(1)
}

func (s *CustomersRepositoryMock) GetByCustomerNumberWithContext(ctx context.Context, customerNumber int) (dto.ResultCustomerRequest, error) {
	args := s.Called(ctx, customerNumber)

	arg0, ok := args.Get(0).(dto.ResultCustomerRequest)
	if !ok {
		return dto.ResultCustomerRequest{}, args.Error(1)
	}

	return arg0, args.Error(1)
}

func (s *CustomersRepositoryMock) GetByNameWithContext(ctx context.Context, name string) (domain.Customer, error) {
	args := s.Called(ctx, name)

	arg0, ok := args.Get(0).(domain.Customer)
	if !ok {
		return domain.Customer{}, args.Error(1)
	}

	return arg0, args.Error(1)
}

func (s *CustomersRepositoryMock) ExistsByCustomerNumberWithContext(ctx context.Context, cid int) bool {
	args := s.Called(ctx, cid)
	return args.Bool(0)
}

func (s *CustomersRepositoryMock) ExistsByIDWithContext(ctx context.Context, id int) bool {
	args := s.Called(ctx, id)
	return args.Bool(0)
}

func (s *CustomersRepositoryMock) ExistsByCustomerNumberAndIDWithContext(ctx context.Context, id, cid int) bool {
	args := s.Called(ctx, id, cid)
	return args.Bool(0)
}

func (s *CustomersRepositoryMock) SaveWithContext(ctx context.Context, sct domain.Customer) (int, error) {
	args := s.Called(ctx, sct)
	return args.Int(0), args.Error(1)
}

func (s *CustomersRepositoryMock) UpdateWithContext(ctx context.Context, sct domain.Customer) error {
	args := s.Called(ctx, sct)
	return args.Error(0)
}

func (s *CustomersRepositoryMock) DeleteWithContext(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
