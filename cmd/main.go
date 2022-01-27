package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MiguelAngelCipamochaF/dynamodb-docker-meli/internal/users"
	"github.com/MiguelAngelCipamochaF/dynamodb-docker-meli/internal/users/models"
	"github.com/MiguelAngelCipamochaF/dynamodb-docker-meli/util"
)

func main() {
	dynamoDB, err := util.InitDynamo()
	if err != nil {
		log.Fatal(err)
	}
	repo := users.NewDynamoRepository(dynamoDB, "Users")
	repo.Store(context.Background(), &models.User{
		Id:         "1",
		Firstname:  "",
		Lastname:   "",
		Username:   "",
		Password:   "",
		Email:      "",
		IP:         "",
		MacAddress: "",
		Website:    "",
		Image:      "",
	})
	user, _ := repo.GetOne(context.Background(), "1")
	fmt.Printf("%+v\n", user)
	repo.Update(context.Background(), "1", "Miguel", "Cipamocha", "correo@gmail.com")
	newUser, _ := repo.GetOne(context.Background(), "1")
	fmt.Printf("%+v\n", newUser)
	repo.Delete(context.Background(), "1")
}
