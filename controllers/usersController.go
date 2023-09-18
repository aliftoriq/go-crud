package controllers

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aliftoriq/go-crud/models"
	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UsersController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Validate(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type usersController struct {
	userRepo repositories.UserRepository
}

func NewUsersController(userRepo repositories.UserRepository) UsersController {
	return &usersController{
		userRepo: userRepo,
	}
}

func (h *usersController) Signup(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FAILED TO READ BODY",
		})
		return
	}

	// Check if the user with the same email already exists
	// userRepo := repositories.NewUserRepository()
	userRepo := h.userRepo
	_, err := userRepo.FindUserByEmail(body.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To Hash Password",
		})
		return
	}

	// Create a new user
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	if err := userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed To create User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User Registered Successfully",
	})
}

func (h *usersController) Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FAILED TO READ BODY",
		})
		return
	}

	userRepo := h.userRepo
	user, err := userRepo.FindByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	tokenString, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	setTokenCookie(c, tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"token":   tokenString,
		"data":    user,
	})
}

// Generate JWT token for the user
func generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	jwtKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func setTokenCookie(c *gin.Context, tokenString string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
}

func (h *usersController) Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
	})

}

func (h *usersController) GetUser(c *gin.Context) {
	userID := c.Param("id")

	// Fetch the user by ID from the repository
	userRepo := h.userRepo
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Return the user data
	c.JSON(http.StatusOK, user)
}

func (h *usersController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	userRepo := h.userRepo
	_, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var updateUser models.User
	if err := c.Bind(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	user.Name = updateUser.Name
	user.Email = updateUser.Email

	if err := userRepo.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *usersController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	userRepo := h.userRepo
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := userRepo.Delete(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// func Signup(c *gin.Context) {
// 	var body struct {
// 		Name     string
// 		Email    string
// 		Password string
// 	}

// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "FAILED TO READ BODY",
// 		})
// 		return
// 	}
// 	var existingUser models.User
// 	if err := initializer.DB.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
// 		c.JSON(http.StatusConflict, gin.H{
// 			"error": "User with this email already exists",
// 		})
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed To Hash Password",
// 		})
// 		return
// 	}

// 	user := models.User{
// 		Name:     body.Name,
// 		Email:    body.Email,
// 		Password: string(hash),
// 	}

// 	result := initializer.DB.Create(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed To create User",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User Registered Succesfully",
// 	})
// }

// func Login(c *gin.Context) {
// 	var body struct {
// 		Email    string
// 		Password string
// 	}

// 	if c.Bind(&body) != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "FAILED TO READ BODY",
// 		})
// 		return
// 	}

// 	var user models.User
// 	initializer.DB.First(&user, "email = ?", body.Email)

// 	if user.ID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid Email or Password",
// 		})
// 		return
// 	}

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid Email or Password",
// 		})
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": user.ID,
// 		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
// 	})

// 	jwtKey := []byte(os.Getenv("SECRET"))
// 	tokenString, err := token.SignedString(jwtKey)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to create token",
// 		})

// 		return
// 	}

// 	c.SetSameSite(http.SameSiteLaxMode)
// 	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Logged in",
// 		"token":   tokenString,
// 		"data":    user,
// 	})

// }

// func UpdateUser(c *gin.Context) {
// 	userID := c.Param("id")

// 	var user models.User
// 	result := initializer.DB.First(&user, userID)

// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "User not found",
// 		})
// 		return
// 	}

// 	var updateUser models.User
// 	if err := c.Bind(&updateUser); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Failed to read request body",
// 		})
// 		return
// 	}

// 	user.Name = updateUser.Name
// 	user.Email = updateUser.Email

// 	result = initializer.DB.Save(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to update user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// func DeleteUser(c *gin.Context) {
// 	userID := c.Param("id")

// 	var user models.User
// 	result := initializer.DB.First(&user, userID)

// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "User not found",
// 		})
// 		return
// 	}

// 	result = initializer.DB.Unscoped().Delete(&user)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to delete user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User deleted successfully",
// 	})
// }
