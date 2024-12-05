package main

import (
	"auth/controllers"
	"auth/pkg/initializers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main(){

	fmt.Println("Hello world!")
	r := gin.Default()
	r.POST("/auth/register", controllers.CreateUser)
	r.Run()
}