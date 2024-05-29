package customer

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/danilosano/web-golang-api/internal/domain"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
)

type Repository interface {
	GetAllWithContext(ctx context.Context) ([]dto.ResultCustomerRequest, error)
	GetWithContext(ctx context.Context, id int) (dto.ResultCustomerRequest, error)
	GetByCustomerNumberWithContext(ctx context.Context, customerNumber int) (dto.ResultCustomerRequest, error)
	ExistsByCustomerNumberWithContext(ctx context.Context, cid int) bool
	ExistsByIDWithContext(ctx context.Context, id int) bool
	ExistsByCustomerNumberAndIDWithContext(ctx context.Context, id, cid int) bool
	SaveWithContext(ctx context.Context, s domain.Customer) (int, error)
	UpdateWithContext(ctx context.Context, s domain.Customer) error
	DeleteWithContext(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllWithContext(ctx context.Context) ([]dto.ResultCustomerRequest, error) {
	query := "SELECT customer_id,customer_number, first_name, last_name, created_at, updated_at FROM customers WHERE deleted_at IS NULL;"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var customers []dto.ResultCustomerRequest

	for rows.Next() {
		c := dto.ResultCustomerRequest{}
		_ = rows.Scan(&c.ID, &c.CustomerNumber, &c.FirstName, &c.LastName, &c.CreatedAt, &c.UpdatedAt)
		customers = append(customers, c)
	}

	return customers, nil
}

func (r *repository) GetWithContext(ctx context.Context, id int) (dto.ResultCustomerRequest, error) {
	query := "SELECT customer_id,customer_number, first_name, last_name, created_at, updated_at FROM customers WHERE deleted_at IS NULL and customer_id=?;"
	row := r.db.QueryRowContext(ctx, query, id)
	c := dto.ResultCustomerRequest{}
	err := row.Scan(&c.ID, &c.CustomerNumber, &c.FirstName, &c.LastName, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return dto.ResultCustomerRequest{}, err
	}

	return c, nil
}

func (r *repository) GetByCustomerNumberWithContext(ctx context.Context, customerNumber int) (dto.ResultCustomerRequest, error) {
	query := "SELECT customer_id,customer_number, first_name, last_name, created_at, updated_at FROM customers WHERE deleted_at IS NULL and customer_number=?;"
	row := r.db.QueryRowContext(ctx, query, customerNumber)
	c := dto.ResultCustomerRequest{}
	err := row.Scan(&c.ID, &c.CustomerNumber, &c.FirstName, &c.LastName, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.ResultCustomerRequest{}, ErrorCustomerNotFound
		}
		return dto.ResultCustomerRequest{}, err
	}

	return c, nil
}

func (r *repository) ExistsByCustomerNumberWithContext(ctx context.Context, cid int) bool {
	query := "SELECT customer_number FROM customers WHERE deleted_at IS NULL and customer_number=?;"
	row := r.db.QueryRowContext(ctx, query, cid)
	err := row.Scan(&cid)
	return errors.Is(err, nil)
}

func (r *repository) ExistsByIDWithContext(ctx context.Context, id int) bool {
	query := "SELECT customer_id FROM customers WHERE deleted_at IS NULL and customer_id=?;"
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&id)
	return errors.Is(err, nil)
}

func (r *repository) ExistsByCustomerNumberAndIDWithContext(ctx context.Context, id, cid int) bool {
	query := "SELECT customer_id FROM customers WHERE deleted_at IS NULL and customer_id=? and customer_number=?;"
	row := r.db.QueryRowContext(ctx, query, id, cid)
	err := row.Scan(&id)
	return errors.Is(err, nil)
}

func (r *repository) SaveWithContext(ctx context.Context, c domain.Customer) (int, error) {
	query := "INSERT INTO customers (customer_number, first_name, last_name, created_at) VALUES (?, ?, ?, ?);"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, &c.CustomerNumber, &c.FirstName, &c.LastName, &c.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) UpdateWithContext(ctx context.Context, c domain.Customer) error {
	query := "UPDATE customers SET customer_number=?, first_name=?, last_name=?, updated_at=? WHERE customer_id=?;"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, &c.CustomerNumber, &c.FirstName, &c.LastName, &c.UpdatedAt, &c.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteWithContext(ctx context.Context, id int) error {
	query := "UPDATE customers SET deleted_at=? WHERE customer_id=?;"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrorCustomerNotFound
	}

	return nil
}
