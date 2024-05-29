package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/danilosano/web-golang-api/internal/customer"
	"github.com/danilosano/web-golang-api/internal/domain/dto"
	"github.com/danilosano/web-golang-api/pkg/web"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service customer.Service
}

func NewCustomerHandler(s customer.Service) *CustomerHandler {
	return &CustomerHandler{
		service: s,
	}
}

// CreateCustomers godoc
// @Summary Create customer
// @Tags Customers
// @Description create customer
// @Accept json
// @Produce json
// @Param customer body dto.CreateCustomerRequest true "Customer to be created"
// @Success 201 {object} web.Responses{data=dto.CreateCustomerRequest} "Success"
// @Failure 400 {object} web.ErrorResponse "Bad Request"
// @Failure 409 {object} web.ErrorResponse "Conflict"
// @Failure 422 {object} web.ErrorResponse "Unprocessable Entity"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
// @Router /api/v1/customers [post]
func (s *CustomerHandler) Store(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err := req.Validate(); err != nil {
		web.Error(c, http.StatusBadRequest, "%v", err)
		return
	}

	sctn, err := s.service.Save(c.Request.Context(), req)

	if err != nil {
		if errors.Is(err, customer.ErrorCustomerNumberAlreadyExist) {
			web.Error(c, http.StatusConflict, err.Error())
			return
		}

		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.Response(c, http.StatusCreated, sctn)
}

// GetCustomers godoc
// @Summary List all customers
// @Description Get all customers
// @Tags Customers
// @Accept json
// @Produce json
// @Success 200 {object} web.Responses{data=[]domain.Customer} "Success"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
// @Router /api/v1/customers [get]
func (s *CustomerHandler) GetAll(c *gin.Context) {
	listCustomers, err := s.service.GetAll(c.Request.Context())
	if err != nil {
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if listCustomers == nil {
		web.Success(c, http.StatusNoContent, listCustomers)
		return
	} else {
		web.Success(c, http.StatusOK, listCustomers)
		return
	}
}

// DeleteCustomer godoc
// @Summary Delete customer
// @Tags Customers
// @Description delete customer
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200
// @Failure 400 {object} web.ErrorResponse "Bad Request"
// @Failure 404 {object} web.ErrorResponse "Not Found"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
// @Router /api/v1/customers/{id} [delete]
func (s *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "invalid ID provided")
		return
	}

	if id == 0 {
		web.Error(c, http.StatusBadRequest, "invalid id provided: id must be a positive non-zero number")
		return
	}

	err = s.service.Delete(c.Request.Context(), int(id))
	if err != nil {
		if errors.Is(err, customer.ErrorCustomerNotFound) {
			web.Error(c, http.StatusNotFound, "customer not found")
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	web.Success(c, http.StatusNoContent, nil)
}

// UpdateCustomer godoc
// @Summary Update customer
// @Tags Customers
// @Description update customerv
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body dto.UpdateCustomerRequest true "Customer to be updated"
// @Success 200 {object} web.Responses{data=dto.UpdateCustomerRequest} "Success"
// @Failure 400 {object} web.ErrorResponse "Bad Request"
// @Failure 404 {object} web.ErrorResponse "Not Found"
// @Failure 409 {object} web.ErrorResponse "Conflict"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
// @Router /api/v1/customers/{id} [patch]
func (s *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "invalid input ID")
		return
	}

	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	if err := req.Validate(); err != nil {
		web.Error(c, http.StatusBadRequest, "%v", err)
		return
	}

	sctn, err := s.service.Update(c.Request.Context(), req, int(id))

	if err != nil {
		if errors.Is(err, customer.ErrorCustomerNumberAlreadyExist) {
			web.Error(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, customer.ErrorCustomerNotFound) {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.Response(c, http.StatusOK, sctn)
}

// GetCustomer godoc
// @Summary Get customer
// @Tags Customers
// @Description Get customer by ID
// @Customer json
// @Param id path int true "Customer ID"
// @Success 200 {object} web.Responses{data=domain.Customer} "Success"
// @Failure 400 {object} web.ErrorResponse "Bad Request"
// @Failure 404 {object} web.ErrorResponse "Not Found"
// @Failure 500 {object} web.ErrorResponse "Internal Server Error"
// @Router /api/v1/customers/{id} [get]
func (s *CustomerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		web.Error(c, http.StatusBadRequest, "invalid input ID")
		return
	}

	if id == 0 {
		web.Error(c, http.StatusBadRequest, "invalid id provided: id must be a positive non-zero number")
		return
	}

	sctn, err := s.service.Get(c.Request.Context(), int(id))
	if err != nil {
		if errors.Is(err, customer.ErrorCustomerNotFound) {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.Success(c, http.StatusOK, sctn)
}
