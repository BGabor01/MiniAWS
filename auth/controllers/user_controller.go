package controllers

import (
	"auth/models"
	"auth/pkg/initializers"
	"auth/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Password string `json:"password" binding:"required,min=6"` 
	Name     string `json:"name" binding:"required,min=2"`  
}

type CreateUserRequest struct {
	UserRequest
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

func GetUser(c *gin.Context){

	var requestData UserRequest

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": []string{err.Error()},
		})
		return
	}

	var user models.User
	result := initializers.DB.Where("name = ?", requestData.Name).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"details": []string{"User not found"},
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"details": []string{"Incorrect password"},
		})
		return
	}

	jwt, err := utils.GenerateJWT(user.Email, user.Name, user.ID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Can't generate JWT",
			"details": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
	})
}