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

type (
	User struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		UpdateAt  string `json:"updated_at"`
	}

	UserResponse struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
		ctx.AbortWithStatusJSON(400, "Couldn't create the new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}

func GetUsers(ctx *gin.Context) {
	var users []UserResponse
	rows, err := db.Db.Query("SELECT id, username, first_name, last_name, email, created_at FROM users")
	if err != nil {
		errMsg := ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Couldn't create the new user: %v", err),
		}
		ctx.JSON(http.StatusBadRequest, errMsg)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user UserResponse
		err = rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
		if err != nil {
			errMsg := ErrorResponse{
				Code:    400,
				Message: fmt.Sprintf("Couldn't create the new user: %v", err),
			}
			ctx.JSON(http.StatusBadRequest, errMsg)
			return
		}
		users = append(users, user)
	}

	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user UserResponse
	err := db.Db.QueryRow("SELECT id, username, first_name, last_name, email, created_at FROM users where id=$1", id).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
	if err != nil {
		errMsg := ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Couldn't get the user: %v", err),
		}
		ctx.JSON(http.StatusBadRequest, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
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

	_, err = db.Db.Exec("update users set username=$1, first_name=$2, last_name=$3, email=$4, updated_at=$5 where id=$6", body.Username, body.FirstName, body.LastName, body.Email, updateAt, id)

	if err != nil {
		errMsg := ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Couldn't update the user: %v", err),
		}
		ctx.JSON(http.StatusBadRequest, errMsg)
		return
	}

	if err != nil {
		errMsg := ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Error checking affected rows: %v", err),
		}
		ctx.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully updated.")
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := db.Db.Exec("delete from users where id=$1", id)

	if err != nil {
		errMsg := ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Couldn't delete the user: %v", err),
		}
		ctx.JSON(http.StatusBadRequest, errMsg)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errMsg := ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Error checking affected rows: %v", err),
		}
		ctx.JSON(http.StatusInternalServerError, errMsg)
		return
	}

	if rowsAffected == 0 {
		errMsg := ErrorResponse{
			Code:    404,
			Message: "User not found or already deleted.",
		}
		ctx.JSON(http.StatusNotFound, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, "User is successfully deleted.")
}
