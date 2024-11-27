// Package core contains the core module of the nutrix application.
//
// The core module contains the main business logic of the application, including
// the data models, services, and handlers for the core features of the
// application.
package core

import (
	"time"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	auth_mw "github.com/elmawardy/nutrix/modules/auth/middlewares"
	"github.com/elmawardy/nutrix/modules/core/handlers"
	"github.com/elmawardy/nutrix/modules/core/middlewares"
	"github.com/elmawardy/nutrix/modules/core/services"
	"github.com/gorilla/mux"
)

// Core is the main struct for the core module.
//
// It contains the necessary fields for the core module to function, including
// the logger, config, settings, prompter, and notification service.
type Core struct {
	// Logger is the logger object for the core module.
	Logger logger.ILogger

	// Config is the config object for the core module.
	Config config.Config

	// Settings is the settings object for the core module.
	Settings config.Settings

	// Prompter is the prompter object for the core module.
	Prompter userio.Prompter

	// NotificationSvc is the notification service object for the core module.
	NotificationSvc services.INotificationService
}

// OnStart is called when the core module is started.
func (c *Core) OnStart() func() {
	return func() {

	}
}

// OnEnd is called when the core module is ended.
func (c *Core) OnEnd() func() {
	return func() {

	}
}

// Seed seeds the database with sample data.
func (c *Core) Seed(entities []string, is_new_only bool) error {

	seedService := services.Seeder{
		Logger:    c.Logger,
		Config:    c.Config,
		Prompter:  c.Prompter,
		IsNewOnly: is_new_only,
	}

	seedablesMap := make(map[string]bool, len(entities))

	for index := range entities {
		seedablesMap[entities[index]] = true
	}

	if _, ok := seedablesMap["materials"]; ok {
		c.Logger.Info("seeding materials ...")
		err := seedService.SeedMaterials(true)
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	if _, ok := seedablesMap["products"]; ok {
		c.Logger.Info("seeding products ...")
		err := seedService.SeedProducts()
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	if _, ok := seedablesMap["categories"]; ok {
		c.Logger.Info("seeding categories ...")
		err := seedService.SeedCategories()
		if err != nil {
			c.Logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// GetSeedables returns a list of seedables.
func (c *Core) GetSeedables() (entities []string, err error) {
	c.Logger.Info("Getting seedables...")

	return []string{
		"products",
		"materials",
		"materialentries",
		"categories",
		"settings",
	}, nil
}

// RegisterBackgroundWorkers registers background workers.
func (c *Core) RegisterBackgroundWorkers() []modules.Worker {

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationSingletonSvc("melody", c.Logger, c.Config)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
		c.NotificationSvc = notification_service
	}

	workers := []modules.Worker{
		{
			Interval: 1 * time.Hour,
			Task: func() {
				services.CheckExpirationDates(c.Logger, c.Config, c.NotificationSvc)
			},
		},
	}

	return workers
}

// RegisterHttpHandlers registers HTTP handlers.
func (c *Core) RegisterHttpHandlers(router *mux.Router, prefix string) {

	auth_svc := auth_mw.NewZitadelAuth(c.Config)

	c.Logger.Info("Successfully conntected to Zitadel")

	router.Handle(prefix+"/api/sales_logs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesLog(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/salesperday", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSalesPerDay(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/entry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteEntry(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialentry", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PushMaterialEntry(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/material", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.AddMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/order", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrder(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/unpaidorders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetUnpaidOrders(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orderpayunpaid", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PayUnpaidOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materiallogs", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterialLogs(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materials", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetMaterials(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/materialcost", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CalculateMaterialCost(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/categories", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetCategories(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/updatecategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/deletecategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteCategory(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/insertcategory", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InsertCategory(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/startorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.StartOrder(c.Config, c.Logger, c.Settings), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orders", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetOrders(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/orderstash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderStash(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/orderremovestash", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.OrderRemoveFromStash(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/ordergetstashed", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetStashedOrders(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/ordercancel", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.CancelOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/submitorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.SubmitOrder(c.Config, c.Logger), "admin", "cashier"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/finishorder", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.FinishOrder(c.Config, c.Logger), "admin", "chef"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/recipeavailability", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeAvailability(c.Config, c.Logger), "admin", "chef", "cashier"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/recipetree", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetRecipeTree(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/products", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProducts(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/addproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.InesrtNewProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/deleteproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.DeleteProduct(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/updateproduct", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.UpdateProduct(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")
	router.Handle(prefix+"/api/productgetready", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetProductReadyNumber(c.Config, c.Logger), "admin", "cashier", "chef"))).Methods("GET", "OPTIONS")
	router.Handle(prefix+"/api/editmaterial", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.EditMaterial(c.Config, c.Logger), "admin"))).Methods("POST", "OPTIONS")

	if c.NotificationSvc == nil {
		notification_service, err := services.SpawnNotificationSingletonSvc("melody", c.Logger, c.Config)
		if err != nil {
			c.Logger.Error(err.Error())
			panic(err)
		}
		c.NotificationSvc = notification_service
	}

	router.Handle(prefix+"/ws", handlers.HandleNotificationsWsRequest(c.Config, c.Logger, c.NotificationSvc))
}