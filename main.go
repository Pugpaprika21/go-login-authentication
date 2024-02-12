package main

import (
	"log"
	"os"
	"time"

	"github.com/Pugpaprika21/go-fiber/app/db"
	"github.com/Pugpaprika21/go-fiber/app/middleware"
	"github.com/Pugpaprika21/go-fiber/app/migrtion"
	"github.com/Pugpaprika21/go-fiber/app/repository"
	"github.com/Pugpaprika21/go-fiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()
	migrtion.Run()
}

func main() {
	setup()

	engine := html.New("./views", ".html")
	engine.Reload(true)
	engine.Debug(true)
	engine.Layout("embed")
	engine.Delims("{{", "}}")
	engine.AddFunc("config", func() fiber.Map {
		return fiber.Map{
			"app_title":     os.Getenv("APP_NAME_WEB"),
			"app_name_html": os.Getenv("APP_NAME_HTML"),
			"base_url":      os.Getenv("APP_URL"),
			"g_site_key":    os.Getenv("G_SITE_KEY"),
			"script_time":   time.Now(),
		}
	})

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       os.Getenv("APP_NAME"),
		Views:         engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(repository.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(middleware.Cors())

	logFile, err := os.OpenFile("logs/fiber.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	app.Use(logger.New(logger.Config{
		Format:        "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		TimeFormat:    "02-Jan-2006",
		TimeZone:      "Asia/Bangkok",
		Output:        logFile,
		DisableColors: true,
	}))

	app.Static("/static", "./public")

	router.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
