package main

import (
	"github.com/Adebayobenjamin/numerisbook/pkg/configs"
	"github.com/Adebayobenjamin/numerisbook/pkg/controllers"
	"github.com/Adebayobenjamin/numerisbook/pkg/repositories"
	"github.com/Adebayobenjamin/numerisbook/pkg/router"
	"github.com/Adebayobenjamin/numerisbook/pkg/services"
	"go.uber.org/dig"
)

var serviceConstructors = []interface{}{
	// DATABASE
	configs.NewMySQLConnection,

	//REDIS
	// configs.SetupRedisConnection,

	//ROUTER
	router.NewApplicationRouter,

	// CONTROLLERS
	controllers.NewInvoiceController,

	// SERVICES
	services.NewAuditService,
	services.NewInvoiceService,
	services.NewReminderService,
	services.NewCustomerService,

	// REPOSITORIES
	repositories.NewAuditTrailRepository,
	repositories.NewInvoiceRepository,
	repositories.NewPaymentRepository,
	repositories.NewReminderRepository,
	repositories.NewCustomerRepository,
	//ENVIRONMENT
	configs.NewEnvironment,

	//LOGGER
	configs.SetupApplicationLogger,

	// SERVER
	NewApplication,
}

func BootstrapServer() (*application, error) {
	c := dig.New()

	for _, v := range serviceConstructors {
		if err := c.Provide(v); err != nil {
			return nil, err
		}
	}
	var server *application
	err := c.Invoke(func(l *application) {
		server = l
	})
	return server, err
}
