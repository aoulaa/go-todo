package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
	"todo/internal/db"
)

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UpdateAt  string `json:"updated_at"`
}

func AddUser(ctx *gin.Context) {
	body := User{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	var updateAt time.Time
	updateAt, err = time.Parse(time.RFC3339, body.UpdateAt)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Bad Input updatedAt time.")
		return
	}

	userID := uuid.New()
	_, err = db.Db.Exec("insert into users(id,username,password, first_name, last_name, email, updated_at) values ($1,$2, $3, $4, $5, $6, $7)", userID, body.Username, body.Password, body.FirstName, body.LastName, body.Email, updateAt)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(400, "Couldn't create the new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}
