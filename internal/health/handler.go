package health

import "github.com/gofiber/fiber/v2"

func SetupHealthHandler(r fiber.Router) {
	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
