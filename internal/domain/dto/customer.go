package dto

import (
	"errors"
	"time"
)

type CreateCustomerRequest struct {
	CustomerNumber *int   `json:"customer_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
}

type UpdateCustomerRequest struct {
	CustomerNumber *int   `json:"customer_number" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
}

type ResultCustomerRequest struct {
	ID             int        `json:"id"`
	CustomerNumber *int       `json:"customer_number"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

func (c *CreateCustomerRequest) Validate() error {
	if c.CustomerNumber == nil {
		return errors.New("invalid input: customer number is required")
	} else if *c.CustomerNumber <= 0 {
		return errors.New("invalid input: customer number must be greater than 0")
	}
	if c.FirstName == "" {
		return errors.New("invalid input: first name is required")
	}
	if c.LastName == "" {
		return errors.New("invalid input: first name is required")
	}
	return nil
}

func (c *UpdateCustomerRequest) Validate() error {
	if c.CustomerNumber == nil {
		return errors.New("invalid input: customer number is required")
	} else if *c.CustomerNumber <= 0 {
		return errors.New("invalid input: customer number must be greater than 0")
	}
	if c.FirstName == "" {
		return errors.New("invalid input: first name is required")
	}
	if c.LastName == "" {
		return errors.New("invalid input: first name is required")
	}
	return nil
}
