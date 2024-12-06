package main

import (
	"cybercampus_module/configs"
	"cybercampus_module/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    routes.UserRoute(app)
    routes.ModuleRoute(app)
    routes.JenisRoleRoute(app)
    routes.TemplateRoute(app)

    
    configs.MongoConnect()

    app.Listen(configs.LoadEnv("PORT"))


}