package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestPanicMiddleware(t *testing.T) {
	fiberApp := fiber.New(fiber.Config{})
	fiberApp.Use(panicMiddleware())
	fiberApp.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNoContent)
	})

	fiberApp.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	response, err := fiberApp.Test(req)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	response.Body.Close()
	require.Equal(t, fiber.StatusNoContent, response.StatusCode)

	req = httptest.NewRequest("GET", "/panic", nil)
	response, err = fiberApp.Test(req)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	response.Body.Close()
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}
