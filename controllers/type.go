package controllers

import "github.com/aliftoriq/go-crud/models"

type (
	Response struct {
		Message string `json:"message"`
	}

	ResponseErr struct {
		Error string `json:"error"`
	}

	User struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
		Data    *User  `json:"data"`
	}

	SignupRequest struct {
		Name     string
		Email    string
		Password string
	}

	LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	CreateArticleResponse struct {
		Message string `json:"message"`
	}

	Article struct {
		Email   string `json:"email"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	GetArticlesResponse struct {
		Message string            `json:"message"`
		Data    *[]models.Article `json:"data"`
	}

	GetArticlesResponseswag struct {
		Message string     `json:"message"`
		Data    *[]Article `json:"data"`
	}

	GetArticleByIDResponse struct {
		Message string          `json:"message"`
		Data    *models.Article `json:"data"`
	}
	GetArticleByIDResponseSwag struct {
		Message string   `json:"message"`
		Data    *Article `json:"data"`
	}
)
