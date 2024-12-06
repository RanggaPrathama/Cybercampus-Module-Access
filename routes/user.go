package routes

import (
	"cybercampus_module/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App){
	//AUTH
	app.Post("/login", controllers.Login)

	
	//USER
	app.Get("/users", controllers.GetAllUsers)
	app.Get("/users/:id", controllers.GetUserById)
	app.Post("/users/add", controllers.CreateUser)
	app.Put("/users/:id", controllers.UpdateUser)
	app.Delete("/users/:id", controllers.DeleteUser)
}