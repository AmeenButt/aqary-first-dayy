package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"tutorial.sqlc.dev/app/tutorial"
)

type User struct {
	conn *pgx.Conn
}
type UserModel struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AuthorStatus bool   `json:"author_status"`
	AuthorID     int64  `json:"author_id"`
}

func CreateUserHanlder(conn *pgx.Conn) *User {
	return &User{conn: conn}
}

func (u *User) CreateUser(c *gin.Context) {
	queries := tutorial.New(u.conn)
	bodyData := &UserModel{}
	if err := c.ShouldBindJSON(bodyData); err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body can not be empty"})
		return
	}
	insertedUser, err := queries.CreateUser(context.Background(), tutorial.CreateUserParams{
		Name:     bodyData.Name,
		Email:    pgtype.Text{String: bodyData.Email, Valid: true},
		Password: pgtype.Text{String: bodyData.Password, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User can not be added"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Added Sucessfully", "result": insertedUser})
}
func (u *User) GetUser(c *gin.Context) {
	queries := tutorial.New(u.conn)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUsers(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User can not be added"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Found Sucessfully", "result": foundUser})
}
