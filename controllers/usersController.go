package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/aliftoriq/go-crud/models"
	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

// @CreateTags godoc

// Signup godoc
// @Summary Register a new user
// @Description Register a new user with a raw JSON request body containing name, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param body body SignupRequest true "User registration details"
// @Success 200 {object} Response{}
// @Failure 400 {object} ResponseErr{}
// @Failure 409 {object} ResponseErr{}
// @Failure 500 {object} ResponseErr{}
// @Router /signup [post]
func (h *usersController) Signup(c *gin.Context) {
	var body SignupRequest

	if c.Bind(&body) != nil {
		resp := ResponseErr{
			Error: "FAILED TO READ BODY",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Check if the user with the same email already exists
	// userRepo := repositories.NewUserRepository()
	userRepo := h.userRepo
	_, err := userRepo.FindUserByEmail(body.Email)
	if err == nil {
		resp := ResponseErr{
			Error: "User with this email already exists",
		}
		c.JSON(http.StatusConflict, resp)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		resp := ResponseErr{
			Error: "Failed To Hash Password",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Create a new user
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	if err := userRepo.CreateUser(&user); err != nil {
		resp := ResponseErr{
			Error: "Failed To create User",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := Response{
		Message: "User Registered Successfully",
	}

	c.JSON(http.StatusOK, resp)
}

// Login godoc
// @Summary Login user
// @Description Log in to the system to get a user token.
// @Tags users
// @Accept json
// @Produce json
// @Param body body LoginRequest true "User login details"
// @Success 200 {object} LoginResponse{}
// @Failure 400 {object} ResponseErr{}
// @Failure 401 {object} ResponseErr{}
// @Router /login [post]
func (h *usersController) Login(c *gin.Context) {
	var body LoginRequest

	if c.BindJSON(&body) != nil {
		resp := ResponseErr{
			Error: "FAILED TO READ BODY",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	userRepo := h.userRepo
	user, err := userRepo.FindByEmail(body.Email)
	if err != nil {
		resp := ResponseErr{
			Error: "Invalid Email or Password",
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		resp := ResponseErr{
			Error: "Invalid Email or Password",
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	tokenString, err := generateToken(user)
	if err != nil {
		resp := ResponseErr{
			Error: "Failed to create token",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	setTokenCookie(c, tokenString)

	loginResp := LoginResponse{
		Message: "Logged in",
		Token:   tokenString,
		Data: &User{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		},
	}

	c.JSON(http.StatusOK, loginResp)
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

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Get User Succesfuly",
	})
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
