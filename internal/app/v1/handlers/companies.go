package handlers

import (
	"context"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CompaniesService interface {
	Create(ctx context.Context, company services.Company) (services.Company, error)
	Get(ctx context.Context, id string) (services.Company, error)
	Update(ctx context.Context, update services.CompanyUpdate) error
	Delete(ctx context.Context, id string) error
}

type EventsPublisher interface {
	OnCreateCompany(e any)
	OnPatchCompany(e any)
	OnDeleteCompany(e any)
}

type companiesHandler struct {
	srv             CompaniesService
	validator       *validator.Validate
	eventsPublisher EventsPublisher
}

func (h companiesHandler) createCompany(c *fiber.Ctx) error {
	var req CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	var registered bool
	if req.Registered != nil {
		registered = *req.Registered
	}
	company, err := h.srv.Create(c.Context(), services.Company{
		Name:              req.Name,
		Description:       req.Description,
		AmountOfEmployees: req.AmountOfEmployees,
		Registered:        registered,
		Type:              req.Type,
	})
	if err != nil {
		return handleError(c, err)
	}

	createdCompany := CompanyFromService(company)
	go h.eventsPublisher.OnCreateCompany(createdCompany)

	return c.JSON(createdCompany)
}

func (h companiesHandler) updateCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.validateId(id); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	var req UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	if err := h.validator.Struct(req); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	update := CompanyUpdateToService(id, req)
	err := h.srv.Update(c.Context(), update)
	if err != nil {
		return handleError(c, err)
	}

	go h.eventsPublisher.OnPatchCompany(update)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h companiesHandler) getCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.validateId(id); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	company, err := h.srv.Get(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(CompanyFromService(company))
}

func (h companiesHandler) deleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.validateId(id); err != nil {
		return handleErrorStatus(c, fiber.StatusBadRequest, err)
	}

	err := h.srv.Delete(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	go h.eventsPublisher.OnDeleteCompany(id)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h companiesHandler) validateId(id string) error {
	return h.validator.Var(id, "required,len=24")
}

func SetupCompaniesRoutes(r fiber.Router, srv CompaniesService, eventsPublisher EventsPublisher) {
	handler := &companiesHandler{
		srv:             srv,
		validator:       validator.New(),
		eventsPublisher: eventsPublisher,
	}

	r.Post("/companies/create", handler.createCompany)
	r.Get("/companies/:id", handler.getCompany)
	r.Patch("/companies/:id", handler.updateCompany)
	r.Delete("/companies/:id", handler.deleteCompany)
}
