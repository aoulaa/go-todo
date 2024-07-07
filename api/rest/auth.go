package rest

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
	"todo/internal/db"
	"todo/internal/pkg/ds"
	"todo/internal/pkg/validation"
)

// Signup is a handler for user signup
func Signup(c *gin.Context) {
	var userInput struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validation.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if validation.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := ds.User{
		Username: userInput.Username,
		Email:    userInput.Email,
		Password: string(hashPassword),
	}
	userID := uuid.New()
	_, err = db.Db.Exec("insert into users(id, username, email, password) values ($1, $2, $3, $4)", userID, user.Username, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
			"msg":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ID":       userID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login is a handler for user login
func Login(c *gin.Context) {

	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindJSON(&userInput) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user ds.User
	err := db.Db.QueryRow("select id, email, password from users where email=$1", userInput.Email).
		Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	fmt.Println(passErr)
	if passErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	tokenID := uuid.New()

	_, err = db.Db.Exec("insert into auth_token(id, user_id, token) values ($1, $2, $3)", tokenID, user.ID, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"successMessage": "Logout successful",
	})
}
