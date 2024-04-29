package filler

import (
	model "github.com/Aminochka4/Golang/final-project/pkg/my-project/model"
)

func PopulateDatabase(models model.Models) error {
	//for _, user := range users {
	//	models.Users.Insert(&user)
	//}
	for _, questionnaire := range questionnaires {
		models.Questionnaires.Insert(&questionnaire)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

//var users = []model.User{
//	{Name: "Example", Surname: "Example", Username: "Example", Email: "Example", Password: &Password{"somepassword",[]byte("$2a$10$Z/OnDzNfHjOVNnoDK1zZ0O0K6U3JmGQw7uvOZIUdPgdZILlC8EGDm")} , Activated: true},
//}

var questionnaires = []model.Questionnaire{
	{Topic: "Example", Questions: "Example", UserId: 1},
	//{Title: "Greek Salad", Description: "Traditional Greek salad with feta cheese", NutritionValue: 200},
	//{Title: "Caprese Salad", Description: "Fresh tomatoes and mozzarella slices", NutritionValue: 180},
}
