package server

import (
	"fmt"
	"time"

	"github.com/alxand/nalo-workspace/internal/api/auth"
	"github.com/alxand/nalo-workspace/internal/api/company"
	"github.com/alxand/nalo-workspace/internal/api/continent"
	"github.com/alxand/nalo-workspace/internal/api/country"
	"github.com/alxand/nalo-workspace/internal/api/dailytask"
	"github.com/alxand/nalo-workspace/internal/api/user"
	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	swagger "github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/alxand/nalo-workspace/docs" // swagger docs
)

// App represents the application server
type App struct {
	app    *fiber.App
	config *config.Config
	logger *zap.Logger
}

// NewApp creates a new application instance
func NewApp(config *config.Config, logger *zap.Logger) *App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
		ErrorHandler: middleware.ErrorHandler(logger),
	})

	return &App{
		app:    app,
		config: config,
		logger: logger,
	}
}

// SetupRoutes configures all application routes
func (a *App) SetupRoutes(
	authHandler *auth.AuthHandler,
	taskHandler *dailytask.TaskHandler,
	userHandler *user.UserHandler,
	continentHandler *continent.ContinentHandler,
	countryHandler *country.CountryHandler,
	companyHandler *company.CompanyHandler,
	authService *auth.Service,
) {
	// Middleware
	a.app.Use(recover.New())
	a.app.Use(helmet.New())
	a.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	a.app.Use(middleware.RequestLogger(a.logger))

	// Health check
	a.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "nalo-workspace-api",
		})
	})

	// Swagger documentation
	a.app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := a.app.Group("/api/v1")

	// Auth routes (no authentication required)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// Protected routes (authentication required)
	protected := api.Group("/", middleware.JWT(authService))

	// Auth protected routes
	protected.Get("/auth/profile", authHandler.Profile)
	protected.Post("/auth/refresh", authHandler.RefreshToken)

	// Daily task routes (authentication required)
	tasksGroup := protected.Group("/dailytask")
	tasksGroup.Post("/", taskHandler.CreateDailyTask)
	tasksGroup.Get("/:date", taskHandler.GetTasksByDate)
	tasksGroup.Put("/:id", taskHandler.UpdateTask)
	tasksGroup.Delete("/:id", taskHandler.DeleteTask)

	// Continent routes (authentication required)
	continentsGroup := protected.Group("/continents")
	continentsGroup.Post("/", continentHandler.CreateContinent)
	continentsGroup.Get("/", continentHandler.GetAllContinents)
	continentsGroup.Get("/:id", continentHandler.GetContinent)
	continentsGroup.Get("/code/:code", continentHandler.GetContinentByCode)
	continentsGroup.Put("/:id", continentHandler.UpdateContinent)
	continentsGroup.Delete("/:id", continentHandler.DeleteContinent)

	// Country routes (authentication required)
	countriesGroup := protected.Group("/countries")
	countriesGroup.Post("/", countryHandler.CreateCountry)
	countriesGroup.Get("/", countryHandler.GetAllCountries)
	countriesGroup.Get("/:id", countryHandler.GetCountry)
	countriesGroup.Get("/code/:code", countryHandler.GetCountryByCode)
	countriesGroup.Get("/continent/:continentId", countryHandler.GetCountriesByContinent)
	countriesGroup.Put("/:id", countryHandler.UpdateCountry)
	countriesGroup.Delete("/:id", countryHandler.DeleteCountry)

	// Company routes (authentication required)
	companiesGroup := protected.Group("/companies")
	companiesGroup.Post("/", companyHandler.CreateCompany)
	companiesGroup.Get("/", companyHandler.GetAllCompanies)
	companiesGroup.Get("/:id", companyHandler.GetCompany)
	companiesGroup.Get("/code/:code", companyHandler.GetCompanyByCode)
	companiesGroup.Get("/country/:countryId", companyHandler.GetCompaniesByCountry)
	companiesGroup.Get("/industry/:industry", companyHandler.GetCompaniesByIndustry)
	companiesGroup.Put("/:id", companyHandler.UpdateCompany)
	companiesGroup.Delete("/:id", companyHandler.DeleteCompany)

	// Admin routes (admin role required)
	adminGroup := protected.Group("/admin", middleware.AdminMiddleware())
	adminGroup.Get("/users", userHandler.ListUsers)
	adminGroup.Get("/users/:id", userHandler.GetUser)
	adminGroup.Put("/users/:id", userHandler.UpdateUser)
	adminGroup.Delete("/users/:id", userHandler.DeleteUser)

	a.logger.Info("Routes configured successfully")
}

// Start starts the server
func (a *App) Start() error {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.logger.Info("Starting server", zap.String("address", addr))
	return a.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (a *App) Shutdown() error {
	a.logger.Info("Shutting down server")
	return a.app.Shutdown()
}
