package version

import "github.com/gofiber/fiber/v2"

const AppName = "go-company-handler"

var (
	BuildTime = "Thu Oct 10 21:28:09 UTC 2024"
	Revision  = "1abcdef"
	Version   = "0.0.0"
	Branch    = "none"
)

type VersionResponse struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
}

func SetupVersionHandler(r fiber.Router) {
	r.Get("/version", handler)
}

func handler(c *fiber.Ctx) error {
	resp := VersionResponse{
		Name:      AppName,
		Version:   Version,
		BuildTime: BuildTime,
		Revision:  Revision,
		Branch:    Branch,
	}

	return c.JSON(resp)
}
