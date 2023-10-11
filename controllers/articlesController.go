package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aliftoriq/go-crud/models"
	"github.com/aliftoriq/go-crud/repositories"
	"github.com/gin-gonic/gin"
)

type ArticlesController interface {
	CreateArticle(c *gin.Context)
	GetArticles(c *gin.Context)
	GetArticleByID(c *gin.Context)
	UpdateArticle(c *gin.Context)
	DeleteArticle(c *gin.Context)
}

type articlesController struct {
	arRepo    repositories.ArticleRepository
	cacheRepo repositories.CacheRepository
}

func NewArticlesController(arRepo repositories.ArticleRepository, cacheRepo repositories.CacheRepository) ArticlesController {
	return &articlesController{
		arRepo:    arRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Create a new article with title and content
// @Tags articles
// @Accept json
// @Produce json
// @Param Authorization header string true "User Token"
// @Param body body Article true "Article creation details"
// @Success 200 {object} Response
// @Failure 400 {object} ResponseErr
// @Failure 500 {object} ResponseErr
// @Router /articles [post]
func (h *articlesController) CreateArticle(c *gin.Context) {
	var body struct {
		Email   string `json:"email"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, ResponseErr{
			Error: "FAILED TO READ BODY",
		})
		return
	}

	article := models.Article{
		Email:   body.Email,
		Title:   body.Title,
		Content: body.Content,
	}

	arRepo := h.arRepo
	if err := arRepo.CreateArticle(article); err != nil {
		c.JSON(http.StatusBadRequest, ResponseErr{
			Error: "Failed to create Article",
		})
		return
	}

	c.JSON(http.StatusOK, CreateArticleResponse{
		Message: "Article Created Successfully",
	})
}

// GetArticles godoc
// @Summary Get a list of articles
// @Description Get a list of articles from the cache or database
// @Tags articles
// @Accept json
// @Produce json
// @Param Authorization header string true "User Token"
// @Success 200 {object} GetArticlesResponseswag
// @Failure 500 {object} ResponseErr
// @Router /articles [get]
func (h *articlesController) GetArticles(c *gin.Context) {
	key := "all_article"

	// Get data from cache redis
	art, status, err := h.cacheRepo.GetValueByKey(c, key)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get articles from cache", err)
		return
	} else if status {
		var cachedArticles []models.Article
		err := json.Unmarshal([]byte(art), &cachedArticles)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to unmarshal cached data", err)
			return
		}

		c.JSON(http.StatusOK, GetArticlesResponse{
			Data:    &cachedArticles,
			Message: "Get Articles Successfully (from cache)",
		})
		return
	}

	// Cache miss, fetch data from the database
	arRepo := h.arRepo
	result, err := arRepo.GetArticles()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get Articles", err)
		return
	}

	// Cache the fetched data
	data, err := json.Marshal(result)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to marshal data for cache", err)
		return
	}

	errSetCache := h.cacheRepo.SetKey(c, key, data, time.Second*60)
	if errSetCache != nil {
		handleError(c, http.StatusInternalServerError, "Failed to set cache", errSetCache)
		return
	}

	c.JSON(http.StatusOK, GetArticlesResponse{
		Data:    result,
		Message: "Get Articles Successfully (from database)",
	})
}

// GetArticleByID godoc
// @Summary Get an article by its ID
// @Description Get an article by providing its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param Authorization header string true "User Token"
// @Success 200 {object} GetArticleByIDResponseSwag
// @Failure 404 {object} ResponseErr
// @Failure 500 {object} ResponseErr
// @Router /articles/{id} [get]
func (h *articlesController) GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	cacheKey := "article_" + id

	art, status, err := h.cacheRepo.GetValueByKey(c, cacheKey)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get article from cache", err)
		return
	} else if status {
		var cachedArticle models.Article
		if err := json.Unmarshal([]byte(art), &cachedArticle); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to unmarshal cached article data", err)
			return
		}

		c.JSON(http.StatusOK, GetArticleByIDResponse{
			Data:    &cachedArticle,
			Message: "Get Article by ID Successfully (from cache)",
		})
		return
	}

	arRepo := h.arRepo
	result, err := arRepo.GetArticleById(id)
	if err != nil {
		handleError(c, http.StatusNotFound, "Article not found", err)
		return
	}

	// Set Cache to Redis
	data, err := json.Marshal(result)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to marshal article data for cache", err)
		return
	}

	if err := h.cacheRepo.SetKey(c, cacheKey, data, time.Second*60); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to set article cache", err)
		return
	}

	c.JSON(http.StatusOK, GetArticleByIDResponse{
		Data:    result,
		Message: "Get Article by ID Successfully (from database)",
	})
}

// UpdateArticle godoc
// @Summary Update article
// @Description Update article with title and content by ID
// @Tags articles
// @Accept json
// @Produce json
// @Param Authorization header string true "User Token"
// @Param id path string true "Article ID"
// @Param body body Article true "Article creation details"
// @Success 200 {object} Response
// @Failure 400 {object} ResponseErr
// @Failure 500 {object} ResponseErr
// @Router /articles/{id} [PUT]
func (h *articlesController) UpdateArticle(c *gin.Context) {
	arRepo := h.arRepo
	id := c.Param("id")

	var updatedArticle models.Article
	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		err := ResponseErr{
			Error: "Invalid request data",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var existingArticle models.Article
	existingArticle.Title = updatedArticle.Title
	existingArticle.Content = updatedArticle.Content

	if err := arRepo.UpdateArticle(id, existingArticle); err != nil {
		err := ResponseErr{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := Response{
		Message: "Article updated successfully",
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteArticle godoc
// @Summary Delete an article by its ID
// @Description Delete an article by providing its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param Authorization header string true "User Token"
// @Success 200 {object} Response
// @Failure 404 {object} ResponseErr
// @Failure 500 {object} ResponseErr
// @Router /articles/{id} [delete]
func (h *articlesController) DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	arRepo := h.arRepo

	if err := arRepo.DeleteArticle(id); err != nil {
		err := ResponseErr{
			Error: err.Error(),
		}
		c.JSON(http.StatusNotFound, err)
		return
	}

	resp := Response{
		Message: "Article deleted successfully",
	}
	c.JSON(http.StatusOK, resp)
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Println(message, err)
	c.JSON(statusCode, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}
