package main

import (
	"auth/models"
	"auth/pkg/initializers"
)


func init(){
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main(){
	initializers.DB.AutoMigrate(&models.User{})
}