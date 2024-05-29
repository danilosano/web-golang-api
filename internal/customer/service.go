package customer

import (
	"context"
	"errors"
	"time"

	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
)

var (
	ErrorCustomerNumberAlreadyExist = errors.New("customer number already exists")
	ErrorCustomerNotFound           = errors.New("customer not found")
)

type Service interface {
	Save(ctx context.Context, s dto.CreateCustomerRequest) (dto.ResultCustomerRequest, error)
	GetAll(ctx context.Context) ([]dto.ResultCustomerRequest, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, s dto.UpdateCustomerRequest, id int) (dto.ResultCustomerRequest, error)
	Get(ctx context.Context, id int) (dto.ResultCustomerRequest, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Save(ctx context.Context, input dto.CreateCustomerRequest) (dto.ResultCustomerRequest, error) {
	sr := domain.Customer{
		ID:             0,
		CustomerNumber: *input.CustomerNumber,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		CreatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),time.Now().Second(), 0, time.Now().Location())}

	if customerNumberExist := s.repository.ExistsByCustomerNumberWithContext(ctx, *input.CustomerNumber); customerNumberExist {
		return dto.ResultCustomerRequest{}, ErrorCustomerNumberAlreadyExist
	}

	customerIdCreated, err := s.repository.SaveWithContext(ctx, sr)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}

	customer, err := s.repository.GetWithContext(ctx, customerIdCreated)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}

	return customer, nil
}

func (s *service) GetAll(ctx context.Context) ([]dto.ResultCustomerRequest, error) {
	return s.repository.GetAllWithContext(ctx)
}

func (s *service) Delete(ctx context.Context, id int) error {
	if customerExist := s.repository.ExistsByIDWithContext(ctx, id); !customerExist {
		return ErrorCustomerNotFound
	}

	return s.repository.DeleteWithContext(ctx, id)
}

func (s *service) Update(ctx context.Context, input dto.UpdateCustomerRequest, id int) (dto.ResultCustomerRequest, error) {
	if customerExist := s.repository.ExistsByIDWithContext(ctx, id); !customerExist {
		return dto.ResultCustomerRequest{}, ErrorCustomerNotFound
	}

	if sameCustomer := s.repository.ExistsByCustomerNumberAndIDWithContext(ctx, id, *input.CustomerNumber); !sameCustomer {
		if customerNumberExist := s.repository.ExistsByCustomerNumberWithContext(ctx, *input.CustomerNumber); customerNumberExist {
			return dto.ResultCustomerRequest{}, ErrorCustomerNumberAlreadyExist
		}
	}

	sr := domain.Customer{
		ID:             id,
		CustomerNumber: *input.CustomerNumber,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		UpdatedAt:      time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),time.Now().Second(), 0, time.Now().Location())}

	err := s.repository.UpdateWithContext(ctx, sr)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}

	sctn, err := s.repository.GetWithContext(ctx, id)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}

	return sctn, nil
}

func (s *service) Get(ctx context.Context, id int) (dto.ResultCustomerRequest, error) {
	if customerExist := s.repository.ExistsByIDWithContext(ctx, id); !customerExist {
		return dto.ResultCustomerRequest{}, ErrorCustomerNotFound
	}

	customer, err := s.repository.GetWithContext(ctx, id)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}
	return customer, nil
}
