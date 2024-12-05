package controllers

import (
	"auth/models"
	"auth/pkg/initializers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2"`  
	Password string `json:"password" binding:"required,min=6"` 
	Email    string `json:"email" binding:"required,email"`
}

func validateRequest(requestData *CreateUserRequest) []string {
	validate := validator.New()
	err := validate.Struct(requestData)

	if err != nil {
		var validationErrors []string
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				validationErrors = append(validationErrors, e.Error())
			}
		} else {
			validationErrors = append(validationErrors, err.Error())
		}
		return validationErrors
	}
	return nil
}


func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CreateUser(c *gin.Context) {
	var requestData CreateUserRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": []string{err.Error()},
		})
		return
	}

	validationErrors := validateRequest(&requestData)
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": validationErrors,
		})
		return
	}

	hashedPassword, err := hashPassword(requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name:     requestData.Name,
		Password: hashedPassword,
		Email:    requestData.Email,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}