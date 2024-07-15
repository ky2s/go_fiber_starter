package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"
	"user_prod/config"
	"user_prod/errs"
	"user_prod/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
)

type App struct {
	*fiber.App
}

func main() {

	config := config.New()

	app := App{
		App: fiber.New(*config.GetFiberConfig()),
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetString("MW_FIBER_CORS_ALLOWORIGINS"),
		AllowMethods:     config.GetString("MW_FIBER_CORS_ALLOWMETHODS"),
		AllowHeaders:     config.GetString("MW_FIBER_CORS_ALLOWHEADERS"),
		MaxAge:           config.GetInt("MW_FIBER_CORS_MAXAGE"),
		AllowCredentials: config.GetBool("MW_FIBER_CORS_ALLOWCREDENTIALS"),
	}))

	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.GetString("BASIC_AUTH_USER"): config.GetString("BASIC_AUTH_PASS"),
		},
		Unauthorized: func(c *fiber.Ctx) error {
			fmt.Println(c)
			return errors.New("error unauthorize")
		},
	}))

	// Initialize Database
	var db *sqlx.DB
	// db, err := config.ConnectDB()
	// if err != nil {
	// 	fmt.Println("failed to init connection", err)
	// 	return
	// }
	// defer db.Close()

	// read
	osArgs := os.Args
	if len(osArgs) > 1 {

		if consoleErr := app.Console(osArgs, db); consoleErr != nil {
			logger.Error(consoleErr.Message)
		}
		return
	}

	// endpoint
	// vendor := handlers.NewVendorVisitHandler(usecases.NewVendorUseCase(repository.NewVendorVisitRepositoryDb(db)))
	// app.Post("/storeVisitation", vendor.StoreVisit)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()

	// Start the server
	err := app.Listen(config.GetString("SERVER_PORT"))
	if err != nil {
		fmt.Println(err.Error())
		app.exit()
	}

}

func (app *App) exit() {
	_ = app.Shutdown()
}

func (app *App) Console(osArgs []string, db *sqlx.DB) *errs.AppError {

	start := time.Now()

	if osArgs[1] == "importVisitUker" {

		if osArgs[2] != "" {

			// rmv := handlers.NewRMVisitHandler(usecases.NewRMUseCase(repository.NewRmVisitRepositoryDb(db)))
			// err := rmv.StoreVisit(fileName)
			// if err != nil {
			// 	return errs.NewUnexpectedError("Error store RM visit : " + err.Message)
			// }

		} else {
			return errs.NewUnexpectedError("File name is required")

		}

		now := time.Now().Local().Format(time.DateTime)
		fmt.Printf("[%s] Insert data Visit Uker executed successfully, no errors found. took: %v\n", now, time.Since(start))
		logger.Info(fmt.Sprintf("[%s] Insert data Visit Uker executed successfully, no errors found. took: %v\n", now, time.Since(start)))

	}

	return nil

}
