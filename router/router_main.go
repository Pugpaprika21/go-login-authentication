package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Pugpaprika21/go-fiber/app/middleware"
	"github.com/Pugpaprika21/go-fiber/app/repository"
	"github.com/Pugpaprika21/go-fiber/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Setup(app *fiber.App) {
	webRouter := app.Group("/web")
	{
		var userJWT middleware.UserAuthenticationMiddlewareInterface = middleware.NewUserAuthentication()
		var userRepository repository.UserRepositoryInterface = repository.NewUserRepository()
		var validationRepository repository.ValidatorRepositoryInterface = repository.NewXValidator()

		pageRouter := webRouter.Group("/page")
		{
			pageRouter.Get("/login", func(c *fiber.Ctx) error { return c.Render("pages/login", fiber.Map{}, "layouts/main") })
			pageRouter.Get("/register", func(c *fiber.Ctx) error { return c.Render("pages/register", fiber.Map{}, "layouts/main") })
			pageRouter.Get("/home", func(c *fiber.Ctx) error { return c.Render("pages/home", fiber.Map{}, "layouts/main") })
		}

		processRouter := webRouter.Group("/process")
		{
			processRouter.Post("/jwt_protected", userJWT.JWTProtected(), func(c *fiber.Ctx) error { return nil })
			processRouter.Post("/login", func(c *fiber.Ctx) error {
				var body dto.UserLoginBodyRequest
				if err := c.BodyParser(&body); err != nil {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
				}

				errs := validationRepository.Validator(&body)
				if len(errs) > 0 {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": errs})
				}

				user, err := userRepository.CheckLogin(body.Username, body.Password)
				if err != nil {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "user incorrect"})
				}

				errHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
				if errHash != nil {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "password incorrect"})
				}

				strJWT, _ := userJWT.SetTokenJWT(user)

				return c.Status(http.StatusOK).JSON(fiber.Map{"user": dto.UserRespone{
					ID:       user.ID,
					Username: body.Username,
					Password: body.Password,
					TokenJWT: strJWT,
				}})
			})

			// เหลือ google recaptcha

			processRouter.Post("/register", func(c *fiber.Ctx) error {
				var body dto.UserRegisterBodyRequest
				if err := c.BodyParser(&body); err != nil {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
				}

				errs := validationRepository.Validator(&body)
				if len(errs) > 0 {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Please fill in the Email correctly."})
				}

				if userRepository.UsernameIsExiting(body.Username) > 0 {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "username is exiting"})
				}

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
				body.Password = string(hashedPassword)
				body.Token = fmt.Sprintf("%s_%s", uuid.New().String(), time.Now().Format("2006-01-02T15:04:05"))

				success, err := userRepository.Save(body)
				if !success {
					return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
				}
				return c.Status(http.StatusOK).JSON(fiber.Map{"success": success})
			})
		}
	}
}
