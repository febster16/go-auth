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
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

// Sign Up		godoc
//
//	@Summary	Signup
//	@Produce	application/json
//	@Param		request	body	requests.UserPayload	true	"User payload"
//	@Tags		users
//	@Router		/signup [post]
func Signup(c *gin.Context) {
	var userPayload requests.UserPayload

	if err := c.Bind(&userPayload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error parsing body payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to hash password"})
		return
	}

	user := models.User{Email: userPayload.Email, Password: string(hashedPassword)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		if result.Error.(*pgconn.PgError).Code == "23505" {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Email already exists"})
			return
		}

		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to create user to DB"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully signed up"})
}

// Log In		godoc
//
//	@Summary	Login
//	@Produce	application/json
//	@Param		request	body	requests.UserPayload	true	"User payload"
//	@Tags		users
//	@Router		/login [post]
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

// Change Password		godoc
//
//	@Summary	ChangePassword
//	@Produce	application/json
//	@Param		request	body	requests.ChangePasswordPayload	true	"Change Password payload"
//	@Tags		users
//	@Router		/change-password [patch]
func ChangePassword(c *gin.Context) {
	var changePasswordPayload requests.ChangePasswordPayload

	if err := c.Bind(&changePasswordPayload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error parsing body payload"})
		return
	}

	var user models.User
	result := database.DB.Find(&user, "email = ?", changePasswordPayload.Email)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordPayload.OldPassword))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordPayload.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to hash password"})
		return
	}

	updateResult := database.DB.Model(models.User{}).Where("email = ?", changePasswordPayload.Email).Updates(models.User{Password: string(hashedPassword)})

	if updateResult.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to update user to DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed password"})

}

// Validate		godoc
//
//	@Summary	Validate
//	@Produce	application/json
//	@Tags		users
//	@Router		/validate [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully logged in as %v", user.(models.User).Email)})
}
