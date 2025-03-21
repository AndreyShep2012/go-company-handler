package app

import (
	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/handlers"
	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"
	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/services"
	"github.com/AndreyShep2012/go-company-handler/internal/config"
	"github.com/AndreyShep2012/go-company-handler/internal/events/simple"
	"github.com/AndreyShep2012/go-company-handler/internal/health"
	"github.com/AndreyShep2012/go-company-handler/internal/version"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRoutes(cfg config.Config, commonRoute, apiRoute fiber.Router, companiesCollection *mongo.Collection) {
	companiesService := services.NewCompaniesService(repositories.NewCompaniesRepository(companiesCollection))
	handlers.SetupCompaniesRoutes(apiRoute, companiesService, simple.New())

	// setup unprotected routes
	version.SetupVersionHandler(commonRoute)
	health.SetupHealthHandler(commonRoute)
}
