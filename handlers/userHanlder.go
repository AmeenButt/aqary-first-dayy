package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"assesment.sqlc.dev/app/models"
	"assesment.sqlc.dev/app/postgres"
	"assesment.sqlc.dev/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	conn *pgx.Conn
}

func CreateUserHanlder(conn *pgx.Conn) *User {
	return &User{conn: conn}
}

func (u *User) CreateUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	data := &models.UserModel{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "Body can not be empty"})
		return
	}
	if data.Name == "" || data.Email == "" || data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}
	_, err := queries.GetUserByEmail(context.Background(), pgtype.Text{String: data.Email, Valid: true})
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exsists with this email"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	insertedUser, err := queries.CreateUser(context.Background(), postgres.CreateUserParams{
		Name:           data.Name,
		Email:          pgtype.Text{String: data.Email, Valid: true},
		Password:       pgtype.Text{String: string(hashedPassword), Valid: true},
		ProfilePicture: pgtype.Text{String: string(data.ProfilePicture), Valid: true},
	})
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User can not be added"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Added Sucessfully", "result": insertedUser})
}

func (u *User) GetUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUserByID(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Found Sucessfully", "result": foundUser})
}
func (u *User) SignIn(c *gin.Context) {
	queries := postgres.New(u.conn)
	data := &models.UserModel{}
	if err := c.ShouldBindJSON(data); err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusNoContent, gin.H{"error": "Body can not be empty"})
		return
	}
	if data.Email == "" || data.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}
	userData, err := queries.GetUserByEmail(context.Background(), pgtype.Text{String: data.Email, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does exsists with this email"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password.String), []byte(data.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Password"})
		return
	}
	token, err := utils.GenerateToken(userData.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Sign in sucessfull", "result": userData, "jwt-token": token})
}
func (u *User) GetAllUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	foundUser, err := queries.ListUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Found Sucessfully", "result": foundUser})
}

func (u *User) UploadProfilePicture(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	queries := postgres.New(u.conn)
	userID := c.Request.FormValue("id")

	defer file.Close()
	currentTime := time.Now().UTC()

	// Format the time in UTC layout
	utcFormat := "2006-01-02T15:04:05.999Z07:00"
	utcTimeString := currentTime.Format(utcFormat)
	replacedString := strings.NewReplacer(".", "_", ",", "_", ";", "_", " ", "_", ":", "_").Replace(utcTimeString)
	// Create a unique filename for the uploaded file
	filename := "uploads/" + replacedString + header.Filename

	out, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer out.Close()

	// Copy the file content to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internat server error"})
		return
	}
	id, err := strconv.Atoi(userID)
	filepath := "uploads/" + utcTimeString + header.Filename
	if id != 0 {
		err := queries.UpdateUserPicture(context.Background(), postgres.UpdateUserPictureParams{
			ID:             int64(id),
			ProfilePicture: pgtype.Text{String: filepath, Valid: true},
		})
		if err != nil {

		}
	}
	c.JSON(200, gin.H{"result": filename, "message": "File Uploaded"})
}
func (u *User) SendOtp(c *gin.Context) {
	queries := postgres.New(u.conn)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUserByID(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate OTP
	otp := utils.GenerateRandomCode()
	result := fmt.Sprintf("Your password reset otp is:%d", otp)

	// Concatenate the OTP to the template

	err = utils.SendMail(foundUser.Email.String, "Your otp for reset password", result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not send otp"})
		return
	}
	updateuser, err := queries.UpdateOTP(context.Background(), postgres.UpdateOTPParams{
		ID:  foundUser.ID,
		Otp: pgtype.Int4{Int32: int32(otp), Valid: true},
	})
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent Sucessfully", "result": updateuser})
}
func (u *User) VerifyOtp(c *gin.Context) {
	queries := postgres.New(u.conn)
	idStr := c.Query("id")
	otpStr := c.Query("otp")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	otp, err := strconv.Atoi(otpStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUserByID(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	convertedOtp := pgtype.Int4{Int32: int32(otp), Valid: true}
	// Generate OTP
	fmt.Println(foundUser.Otp)
	fmt.Println(convertedOtp)
	if foundUser.Otp != convertedOtp {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incorrect OTP"})
		return
	}
	updateuser, err := queries.UpdateOTP(context.Background(), postgres.UpdateOTPParams{
		ID:  foundUser.ID,
		Otp: pgtype.Int4{Int32: 0, Valid: true},
	})
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified", "result": updateuser})
}
