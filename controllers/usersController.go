package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thisgirlElan/jwt_auth/initializers"
	"github.com/thisgirlElan/jwt_auth/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	// fetch email

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// fetch password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password.",
		})

		return
	}

	// create user

	user := models.User{Email: body.Email, Password: string(hash)}

	result := initializers.Db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User exists",
		})

		return
	}

	// success response

	c.JSON(http.StatusOK, gin.H{})
}

func LoginHandler(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	// get email and pass from req body

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch request. Please try again",
		})

		return
	}

	// find user
	var user models.User
	initializers.Db.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// compare passwords

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETJTOKEN")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})

		return
	}

	// send token as cookie
	c.SetSameSite(http.SameSiteLaxMode)
	// secure should be true for non local host
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func UsersHandler(c *gin.Context) {
	var users []models.User
	initializers.Db.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func UserHandler(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func UserByIdHandler(c *gin.Context) {
	var body struct {
		ID int
	}

	// get ID req body

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch request. Please try again",
		})

		return
	}

	// find user
	var user models.User
	initializers.Db.First(&user, "id = ?", body.ID)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func UpdateUserHandler(c *gin.Context) {

	id := c.Param("id")

	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch request. Please try again",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password.",
		})

		return
	}

	var user models.User
	initializers.Db.First(&user, id)

	initializers.Db.Model(&user).Updates(models.User{
		Email:    body.Email,
		Password: string(hash),
	})

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})

}

func DeleteProfile(c *gin.Context) {
	id := c.Param("id")

	initializers.Db.Delete(&models.User{}, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted.",
	})
}
