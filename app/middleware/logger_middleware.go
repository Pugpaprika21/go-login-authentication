package middleware

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		file, err := os.OpenFile("logs/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Middleware Logger Error opening log file: %v", err)
		}

		log.SetOutput(io.MultiWriter(os.Stdout, file))
		defer file.Close()

		logText := fmt.Sprintf(`
IP: %s
METHOD: %s
PATH: %s
STATUS_CODE: %d
QUERY_STRING: %s
PARAMS: %s
REQUEST_BODY: %s

`,
			ctx.IP(),
			ctx.Method(),
			ctx.Path(),
			ctx.Response().StatusCode(),
			ctx.OriginalURL(),
			ctx.Params("*"),
			string(ctx.Request().Body()),
		)

		log.Print(logText)
		return ctx.Next()
	}
}
