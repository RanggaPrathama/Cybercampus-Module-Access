package helpers

import (
	"context"
	"cybercampus_module/configs"
	"cybercampus_module/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionUserModule = configs.GetCOllection(configs.Client, "user_module")
var collectionTemplate = configs.GetCOllection(configs.Client, "templates")


func SyncModuleTemplate(jenis_user primitive.ObjectID, idUser primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Syncing module for user with id:", idUser)
	fmt.Println("Jenis user:", jenis_user)

	
	var userModule models.UserModule
	err := collectionUserModule.FindOne(ctx, bson.M{"id_user": idUser}).Decode(&userModule)

	if err != nil {
	
		if err.Error() == "mongo: no documents in result" {
			
			var template models.TemplateUserModuleRequest
			err = collectionTemplate.FindOne(ctx, bson.M{"_id": jenis_user}).Decode(&template)

			fmt.Println("Template:", template)

			if err != nil {
				return false, err
			}

			
			newModules := []primitive.ObjectID{}
			newModules = append(newModules, template.Template...)
			//for _, module := range template.Template {

				// data, _ := bson.Marshal(module)
				// var module models.ModuleRequest
				// bson.Unmarshal(data, &module)
			//	newModules = append(newModules, module.ID)
			//}

	
			newUserModule := models.UserModule{
				ID:         primitive.NewObjectID(),
				IDUser:     idUser,
				JenisUser:  template.JenisUser,
				MODULES:    newModules,
				CREATED_AT: time.Now(),
				UPDATED_AT: time.Now(),
			}

			_, err = collectionUserModule.InsertOne(ctx, newUserModule)
			if err != nil {
				return false, err
			}

			fmt.Println("User module created : ", newUserModule)
			return true, nil
		}

		// Jika error lain, return error
		return false, err
	}

	// Jika entri sudah ada, update modules dengan mengambil dari template
	var template models.TemplateResponse
	err = collectionTemplate.FindOne(ctx, bson.M{"_id": jenis_user}).Decode(&template)

	if err != nil {
		return false, err
	}

	// Siapkan modul baru dari template
	updatedModules := []primitive.ObjectID{}
	for _, module := range template.Template {
		updatedModules = append(updatedModules, module.ID)
	}

	// Update entri di user_module
	update := bson.M{
		"$set": bson.M{
			"modules":    updatedModules,
			"updated_at": time.Now(),
		},
	}

	_, err = collectionUserModule.UpdateOne(ctx, bson.M{"id_user": idUser}, update)
	if err != nil {
		return false, err
	}

	return true, nil
}