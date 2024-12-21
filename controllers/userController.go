package controllers

import (
	"context"
	"cybercampus_module/configs"
	"cybercampus_module/helpers"
	"cybercampus_module/models"
	"cybercampus_module/response"
	"fmt"

	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = configs.GetCOllection(configs.Client, "users")

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var users []models.UserResponse

	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when fetching users",
			Data:    err.Error(),
		})
	}

	for cursor.Next(ctx) {
		var value models.UserRequest
		if err := cursor.Decode(&value); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error when decoding user",
				Data:    err.Error(),
			})
		}

		var jenisUser models.JenisUserResponse

		//check jenis user
		err := collectionTemplate.FindOne(ctx, bson.M{"_id": value.JENIS_USER}).Decode(&jenisUser)
		if err != nil && jenisUser.JenisUser != "" {
			return c.Status(fiber.StatusNotFound).JSON(response.Response{
				Status:  fiber.StatusNotFound,
				Message: "JENIS USER NOT FOUND",
				Data:    err.Error(),
			})
		}

		user := models.UserResponse{
			ID:       value.ID,
			Username: value.Username,
			NM_USER:  value.NM_USER,
			Email:    value.Email,
			Password: value.Password,
			Role:     value.Role,
			JENIS_USER: func() string {
				if jenisUser.JenisUser == "" {
					return "admin"
				}
				return jenisUser.JenisUser
			}(),
			LastLogin: value.LastLogin,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
			Phone:     value.Phone,
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    users,
	})

}

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Error when parsing body",
			Data:    err.Error(),
		})
	}

	hashedPassword := helpers.HashPasword(user.Password)

	//cekJenisUser 
	var jenisUser models.TemplateRequest
	
	err := collectionTemplate.FindOne(ctx, bson.M{"_id": user.JENIS_USER}).Decode(&jenisUser)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Status:  fiber.StatusNotFound,
			Message: "Not Found Jenis User With Templates",
			Data:    err.Error(),
		})
	}

	module := []primitive.ObjectID{}
	module = append(module, jenisUser.Template...)
	
	userNew := models.UserRequest{
		Username:    user.Username,
		NM_USER:     user.NM_USER,
		Password:    hashedPassword,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
		Photo:       user.Photo,
		Phone:       user.Phone,
		JENIS_USER:  user.JENIS_USER,
		Address:     user.Address,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		Modules:   module,
	}

	data, err := collection.InsertOne(ctx, userNew)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when inserting user",
			Data:    err.Error(),
		})
	}


	return c.Status(fiber.StatusCreated).JSON(response.Response{
		Status:  fiber.StatusCreated,
		Message: "Success",
		Data:    data,
	})
}

func GetUserById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Locals("id_user").(string)
	jenis_user := c.Locals("jenis_user").(string)

	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when parsing id",
			Data:    err.Error(),
		})
	}

	var user models.ProfileResponse
	var jenisUser models.JenisUserResponse

	//checkJenisUser
	err = collectionTemplate.FindOne(ctx, bson.M{"jenis_user": jenis_user}).Decode(&jenisUser)

	if err != nil && jenisUser.JenisUser != "" {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Status:  fiber.StatusNotFound,
			Message: "Not Found Jenis User",
			Data:    err.Error(),
		})
	}

	err = collection.FindOne(ctx, bson.M{"_id": hexId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when fetching user",
			Data:    err.Error(),
		})
	}

	profileUser := models.ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		NM_USER:  user.NM_USER,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		JENIS_USER: func() string {
			if jenisUser.JenisUser == "" {
				return "admin"
			}
			return jenisUser.JenisUser
		}(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    profileUser,
	})
}

func UpdateUser(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	hexId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when parsing id",
			Data:    err.Error(),
		})
	}

	var user models.UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Error when parsing body",
			Data:    err.Error(),
		})
	}

	fmt.Printf("Parsed User: %+v\n", user)

	updateUser := bson.M{}

	value := reflect.ValueOf(user)
	tipe := reflect.TypeOf(user)

	for i := 0; i < value.NumField(); i++ {
		fieldName := tipe.Field(i).Tag.Get("bson")
		//fmt.Print("Tes : ", fieldName)
		if fieldName == "" || fieldName == "-" {
			continue
		}
		fieldValue := value.Field(i).Interface()
		switch val := fieldValue.(type) {
		case string:
			if val != "" {
				updateUser[fieldName] = val

			}
		case time.Time:
			if !val.IsZero() {
				updateUser[fieldName] = val
			}
		case bool:
			if val {
				updateUser[fieldName] = val
			}
		default:
			if !reflect.ValueOf(fieldValue).IsZero() {
				updateUser[fieldName] = fieldValue
			}
		}

	}

	//fmt.Printf("Update User Fields: %+v\n", updateUser)
	//fmt.Println("Len User", len(updateUser))

	if len(updateUser) > 0 {
		updateUser["updated_at"] = time.Now()
		//fmt.Printf("Update User Fields: %+v\n", updateUser)
		result, err := collection.UpdateOne(ctx, bson.M{"_id": hexId}, bson.M{"$set": updateUser})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Status:  fiber.StatusInternalServerError,
				Message: "Error when updating user",
				Data:    err.Error(),
			})
		}

		if result.ModifiedCount == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(response.Response{
				Status:  fiber.StatusBadRequest,
				Message: "No document was updated",
				Data:    nil,
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "No field to update",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})

}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	hexId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when parsing id",
			Data:    err.Error(),
		})
	}

	respon, err := collection.DeleteOne(ctx, bson.M{"_id": hexId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when deleting user",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    respon,
	})
}

func UpdateJenisUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Error when parsing body",
			Data:    err.Error(),
		})
	}

	updateUser := bson.M{
		"$set": bson.M{
			"jenis_user": user.JENIS_USER,
		},
	}

	filter := bson.M{"_id": user.ID}

	result, err := collection.UpdateOne(ctx, filter, updateUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when updating user",
			Data:    err.Error(),
		})
	}

	responses, err := helpers.SyncModuleTemplate(user.JENIS_USER, user.ID)

	if err != nil || !responses {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Error when syncing user module",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Response{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})

}
