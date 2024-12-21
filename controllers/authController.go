package controllers

import (
	"context"
	"cybercampus_module/helpers"
	"cybercampus_module/models"
	"cybercampus_module/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request",
			"data":    err.Error(),
		})
	}

	var LoginData models.UserRequest

	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&LoginData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "User Not Found",
			"data":    err.Error(),
		})
	}

	if !helpers.ComparePassword(LoginData.Password, user.Password) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Password is incorrect",
			"data":    nil,
		})
	}

	var jenisUser models.JenisUserResponse
	_ = collectionTemplate.FindOne(ctx, bson.M{"_id": LoginData.JENIS_USER}).Decode(&jenisUser)
	token, err := helpers.GenerateToken(LoginData.ID.Hex(), LoginData.Username, LoginData.Email, jenisUser.JenisUser, LoginData.Role)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Error when generating token",
			"data":    err.Error(),
		})
	}


	pipeline := mongo.Pipeline{
		{{
			Key: "$match", Value: bson.D{
				{Key: "_id", Value: LoginData.ID},
			},
		}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "module"},
				{Key: "localField", Value: "modules"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "module_data"},
			},
		}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "templates"},
				{Key: "localField", Value: "jenis_user"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "template_data"},
			},
		}},
		//  {{
		// 	Key: "$unwind", Value: "$module_data",
		//   }},
		{{Key: "$unwind", Value: "$template_data"}},
		{{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "username", Value: 1},
				{Key: "nm_user", Value: 1},
				{Key: "email", Value: 1},
				{Key: "password", Value: 1},
				{Key: "jenis_user", Value: "$template_data.jenis_user"},
				{Key: "role", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "address", Value: 1},
				{Key: "date_of_birth", Value: 1},
				{Key: "modules", Value: "$module_data"},
			},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Error when fetching aggregate user module",
			"data":    err.Error(),
		})
	}

	var LoginDataResponse []models.UserLogin

	if err := cursor.All(ctx, &LoginDataResponse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Error when decoding user module",
			"data":    err.Error(),
		})
	}

	if LoginData.Role == "admin" {
		var modules []models.ModuleResponse

		cursor, err := collectionModule.Find(ctx, bson.D{})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error when fetching modules",
				Data:    err.Error(),
			})
		}

		if err = cursor.All(ctx, &modules); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error when decoding modules",
				Data:    err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.Response{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: models.UserLogin{
				ID:       LoginData.ID.Hex(),
				Username: LoginData.Username,
				NM_USER:  LoginData.NM_USER,
				Email:    LoginData.Email,
				Password: LoginData.Password,
				JENIS_USER: func() string {
					if jenisUser.JenisUser == "" {
						return "admin"
					}
					return jenisUser.JenisUser
				}(),
				IsActive:    LoginData.IsActive,
				Role:        LoginData.Role,
				Phone:       LoginData.Phone,
				Address:     LoginData.Address,
				DateOfBirth: LoginData.DateOfBirth,
				Modules:     modules,
				TOKEN:       token,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.LoginResponse{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data: LoginDataResponse[0],
		Token: token,
	},
	)

}
