package domain

import "time"

type Customer struct {
	ID             int       `json:"id" db:"customer_id"`
	CustomerNumber int       `json:"customer_number"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
}
