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
		c.JSON(http.StatusCreated, models.UserReturn{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
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

// UpdateUserRoleHandler updates the role of a user by username
//
//	@summary		Assign a role to a user
//	@description	Sets the role of an existing user. Valid roles: administrator, preparator, reception.
//	@tags			Users
//	@accept			json
//	@produce		json
//	@param			username	path		string						true	"Username"
//	@param			role		body		models.UserRoleUpdateForm	true	"Role to assign"
//	@success		200			{object}	models.UserReturn
//	@failure		400			{object}	models.ErrorResponse	"Invalid role"
//	@failure		404			{object}	models.ErrorResponse	"Username not found"
//	@failure		500			{object}	models.ErrorResponse	"Internal error"
//	@router			/users/{username}/role [put]
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
//
//	@summary		Register a new user
//	@description	Creates a user account. New users have no role and cannot log in until one is assigned.
//	@tags			Users
//	@accept			json
//	@produce		json
//	@param			user	body		models.UserForm	true	"Username and password"
//	@success		201		{object}	models.UserReturn
//	@failure		400		{object}	models.ErrorResponse	"Invalid payload"
//	@failure		500		{object}	models.ErrorResponse	"Error creating user"
//	@router			/users/register [post]
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// User Form is username + password
		var SignUp models.UserForm
		if err := c.ShouldBindJSON(&SignUp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		// Check that the user does not exist
		_, err := GetUserByUsername(db, SignUp.Username)
		// Err must NOT be nil, other
		if err == nil {
			// If we find the user, err is nil, meaning we have a 409 conflict
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		// Ensure the error is not RecordNotFound, since we don't expect to find users who try to register.
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		c.JSON(http.StatusCreated, models.UserReturn{
			ID:        User.ID,
			Username:  User.Username,
			CreatedAt: User.CreatedAt,
			UpdatedAt: User.UpdatedAt,
		})
	}
}

// Login will ensure the users are authenticated and return a JWT to use.
//
//	@summary		Login
//	@description	Authenticates a user and returns a JWT valid for 8 hours. The token must be sent as a Bearer token on protected endpoints.
//	@tags			Users
//	@accept			json
//	@produce		json
//	@param			credentials	body		models.UserForm	true	"Username and password"
//	@success		200			{object}	models.TokenResponse
//	@failure		400			{object}	models.ErrorResponse	"Invalid credentials"
//	@failure		403			{object}	models.ErrorResponse	"User has no role assigned"
//	@failure		404			{object}	models.ErrorResponse	"Username not found"
//	@failure		500			{object}	models.ErrorResponse	"Error generating token"
//	@router			/users/login [post]
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
