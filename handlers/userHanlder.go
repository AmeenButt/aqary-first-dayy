package handlers

import (
	"context"
	"fmt"
	"io"
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

type UserHanlder struct {
	conn *pgx.Conn
	ctx  *context.Context
}

func CreateUserHanlder(conn *pgx.Conn, ctx *context.Context) *UserHanlder {
	return &UserHanlder{conn: conn, ctx: ctx}
}

func (u *UserHanlder) CreateUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	data := &models.CreateUserModel{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	_, err := queries.GetUserByEmail(*u.ctx, pgtype.Text{String: data.Email, Valid: true})
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exsists with this email"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error While encrypting password"})
		return
	}
	insertedUser, err := queries.CreateUser(*u.ctx, postgres.CreateUserParams{
		Name:           data.Name,
		Email:          pgtype.Text{String: data.Email, Valid: true},
		Password:       pgtype.Text{String: string(hashedPassword), Valid: true},
		ProfilePicture: pgtype.Text{String: string(data.ProfilePicture), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Added Sucessfully", "result": insertedUser})
}

func (u *UserHanlder) GetUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUserByID(*u.ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Found Sucessfully", "result": foundUser})
}

func (u *UserHanlder) SignIn(c *gin.Context) {
	queries := postgres.New(u.conn)
	data := &models.SignInUserModel{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.GetBindErrorMessage(err)})
		return
	}
	userData, err := queries.GetUserByEmail(*u.ctx, pgtype.Text{String: data.Email, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password.String), []byte(data.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Password"})
		return
	}
	token, err := utils.GenerateToken(userData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating jwt token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sign in sucessfull", "result": userData, "jwt-token": token})
}

func (u *UserHanlder) GetAllUser(c *gin.Context) {
	queries := postgres.New(u.conn)
	foundUser, err := queries.ListUsers(*u.ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Found Sucessfully", "result": foundUser})
}

func (u *UserHanlder) UploadProfilePicture(c *gin.Context) {
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
		c.JSON(500, gin.H{"error": "Error while uploading file"})
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
	if err != nil {
		c.JSON(500, gin.H{"error": "Can not use this userId"})
		return
	}
	filepath := "uploads/" + utcTimeString + header.Filename
	if id != 0 {
		_ = queries.UpdateUserPicture(*u.ctx, postgres.UpdateUserPictureParams{
			ID:             int64(id),
			ProfilePicture: pgtype.Text{String: filepath, Valid: true},
		})
	}
	c.JSON(200, gin.H{"result": filename, "message": "File Uploaded"})
}

func (u *UserHanlder) SendOtp(c *gin.Context) {
	queries := postgres.New(u.conn)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		return
	}
	foundUser, err := queries.GetUserByID(*u.ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
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
	updateuser, err := queries.UpdateOTP(*u.ctx, postgres.UpdateOTPParams{
		ID:  foundUser.ID,
		Otp: pgtype.Int4{Int32: int32(otp), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent Sucessfully", "result": updateuser})
}

func (u *UserHanlder) VerifyOtp(c *gin.Context) {
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
	foundUser, err := queries.GetUserByID(*u.ctx, int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
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
	updateuser, err := queries.UpdateOTP(*u.ctx, postgres.UpdateOTPParams{
		ID:  foundUser.ID,
		Otp: pgtype.Int4{Int32: 0, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified", "result": updateuser})
}

