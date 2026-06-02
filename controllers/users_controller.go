package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"wacdo-backend/middlewares"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid User provided"})
			return
		}
		if err := CreateUser(db, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func UpdateUserRole(db *gorm.DB, username string, role models.Role) (*models.User, error) {
	var user *models.User
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}
	// Update the user, return the error
	return user, db.Model(&user).Update("role", role).Error
}

func UpdateUserRoleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var updateRole models.UserRoleUpdateForm
		if err := c.ShouldBindJSON(&updateRole); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		if !updateRole.Role.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect role received"})
			return
		}
		user, err := UpdateUserRole(db, username, *updateRole.Role)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Username %s not found", username)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return a view model of the users, to avoid exposing sensitive data
		c.JSON(http.StatusOK, models.UserReturn{
			Username:  user.Username,
			Role:      *user.Role,
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// Register allows users to sign up using username and password. They have no role by default.
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// User Form is username + password
		var SignUp models.UserForm
		if err := c.ShouldBindJSON(&SignUp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		// Check that the user does not exist
		user, err := GetUserByUsername(db, SignUp.Username)
		if user.ID != uuid.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
			return
		}
		// Create the user model
		var User models.User
		User.Username = SignUp.Username
		// Hash the password
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(SignUp.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
			return
		}
		User.Password = string(encryptedPassword)
		if err := CreateUser(db, &User); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
			return
		}
		c.JSON(http.StatusCreated, User)
	}
}

// Login will ensure the users are authenticated and return the a JWT to use.
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// User form is username + password
		var form models.UserForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username/password received"})
			return
		}
		user, err := GetUserByUsername(db, form.Username)
		if err != nil {
			// User not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
				return
			}
			// Other errors
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username/password received"})
			return
		}
		// Ensure the password is correct
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username/password received"})
			return
		}
		// Ensure the user has a role, if not, we deny the access.
		if user.Role == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "users without a role cannot login."})
			return
		}
		// Return a valid JWT with the correct information to the user
		var jwtClaims middlewares.JwtStruct = middlewares.JwtStruct{
			Role: *user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "wacdo-backend",
				Subject:   user.Username,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
		ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating JWT"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": ss})
	}
}
