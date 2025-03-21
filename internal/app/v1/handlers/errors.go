package handlers

import (
	"errors"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"
	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError

	switch {
	case errors.As(err, &services.ErrNotFound{}):
		status = fiber.StatusNotFound
	case errors.As(err, &services.ErrDbDuplicatedKey{}):
		status = fiber.StatusConflict
	}

	return handleErrorStatus(c, status, err)
}

func handleErrorStatus(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}
