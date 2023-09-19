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

func (h *articlesController) CreateArticle(c *gin.Context) {
	var body struct {
		Email   string `gorm:"unique"`
		Title   string
		Content string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FAILED TO READ BODY",
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create Article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article Created Succesfuly",
	})
}

func (h *articlesController) GetArticles(c *gin.Context) {
	key := "all_article"

	// Get data from cache redis
	art, status, err := h.cacheRepo.GetValueByKey(c, key)
	if err != nil {
		log.Println("Get Cache Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get articles from cache",
			"details": err.Error(),
		})
		return
	} else if status {
		var cachedArticles []models.Article
		err := json.Unmarshal([]byte(art), &cachedArticles)
		if err != nil {
			log.Println("Unmarshal Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to unmarshal cached data",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    cachedArticles,
			"message": "Get Articles Successfully (from cache)",
		})
		return
	}

	// Cache miss, fetch data from the database
	arRepo := h.arRepo
	result, err := arRepo.GetArticles()
	// initializer.DB.Find(&articles)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get Articles",
			"details": err.Error(),
		})
		return
	}

	// Cache the fetched data
	data, err := json.Marshal(result)
	if err != nil {
		log.Println("Marshal Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to marshal data for cache",
			"details": err.Error(),
		})
		return
	}

	errSetCache := h.cacheRepo.SetKey(c, key, data, time.Second*60)

	if errSetCache != nil {
		log.Println("Set Cache Error:", errSetCache)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to set cache",
			"details": errSetCache,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "Get Articles Successfully (from database)",
	})
}

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

		c.JSON(http.StatusOK, gin.H{
			"data":    cachedArticle,
			"message": "Get Article by ID Successfully (from cache)",
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

	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "Get Article by ID Successfully (from database)",
	})
}

func (h *articlesController) UpdateArticle(c *gin.Context) {
	arRepo := h.arRepo
	id := c.Param("id")

	var updatedArticle models.Article
	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var existingArticle models.Article
	existingArticle.Title = updatedArticle.Title
	existingArticle.Content = updatedArticle.Content

	if err := arRepo.UpdateArticle(id, existingArticle); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
	})
}

func (h *articlesController) DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	arRepo := h.arRepo

	if err := arRepo.DeleteArticle(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted successfully",
	})
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	log.Println(message, err)
	c.JSON(statusCode, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}
