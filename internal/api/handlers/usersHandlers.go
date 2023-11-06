package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/febster16/go-auth/config"
	"github.com/febster16/go-auth/database"
	"github.com/febster16/go-auth/internal/api/requests"
	"github.com/febster16/go-auth/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var userPayload requests.UserPayload

	if err := c.Bind(&userPayload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error parsing body payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), 10)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to hash password"})
		return
	}

	user := models.User{Email: userPayload.Email, Password: string(hashedPassword)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create user to DB"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully created user"})
}

func Login(c *gin.Context) {
	var userPayload requests.UserPayload

	if err := c.Bind(&userPayload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error parsing body payload"})
		return
	}

	var user models.User
	result := database.DB.Find(&user, "email = ?", userPayload.Email)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.CONFIG.SECRET))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully logged in as %v", user.(models.User).Email)})
}