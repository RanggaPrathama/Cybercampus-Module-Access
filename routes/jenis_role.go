package routes

import (
	"cybercampus_module/controllers"

	"github.com/gofiber/fiber/v2"
)

func JenisRoleRoute(app *fiber.App) {
	app.Get("/jenis_roles", controllers.GetAllJenisUser)
	app.Post("/jenis_roles/add", controllers.CreateJenisUser)
}