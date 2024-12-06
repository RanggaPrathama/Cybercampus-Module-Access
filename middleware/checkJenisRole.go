package middleware

import (
	//"context"
	//"cybercampus_module/configs"
	//"cybercampus_module/models"
	//"time"

	"github.com/gofiber/fiber/v2"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
)


//var collectionTemplate *mongo.Collection = configs.GetCOllection(configs.Client, "template")

func CheckJenisRole(role []string) fiber.Handler{
	return func(c *fiber.Ctx) error {

		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()

		//jenis_user := c.Locals("jenis_user").(string)
		roleClaims := c.Locals("role").(string)

		// if jenis_user == "" {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"status":  fiber.StatusBadRequest,
		// 		"message": "Jenis User does not exist",
		// 		"data":    nil,
		// 	})
		// }

		//var checkJenisUser models.JenisUserResponse

		//_ := collectionTemplate.FindOne(ctx, bson.M{"jenis_user": jenis_user}).Decode(&checkJenisUser)

		// if err != nil {
		// 	if err == mongo.ErrNoDocuments {
		// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 			"status":  fiber.StatusBadRequest,
		// 			"message": "Jenis User does not exist",
		// 			"data":    nil,
		// 		})
		// 	}

		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"status":  fiber.StatusInternalServerError,
		// 		"message": "Error when checking Jenis User",
		// 		"data":    err.Error(),
		// 	})
		// }

		isAllowed := false 
		for _, v :=range role {
			if v == roleClaims {
				isAllowed = true
				break
			}

			if !isAllowed {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  fiber.StatusUnauthorized,
					"message": "You are not allowed to access this resource",
					"data":    nil,
				})
			}
		}

		return c.Next()
	}
}