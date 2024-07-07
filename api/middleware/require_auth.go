package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
	"todo/internal/db"
	"todo/internal/pkg/ds"
)

type AuthUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"Email"`
}

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		userId := claims["sub"]
		var user ds.User

		userErr := db.Db.QueryRow("select id, username, email from users where id=$1", userId).Scan(&user.ID, &user.Username, &user.Email)
		if userErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		if len(user.ID) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		var isActive bool
		tokenErr := db.Db.QueryRow("select is_active from auth_token where token=$1", token.Raw).Scan(&isActive)
		if tokenErr != nil {
			fmt.Println(tokenErr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized invalid token",
			})
			return
		}

		if !isActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized token is not active",
			})
			return
		}

		authUser := AuthUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}

		c.Set("authUser", authUser)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
