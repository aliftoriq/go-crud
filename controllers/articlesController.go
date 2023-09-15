package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aliftoriq/go-crud/cache"
	"github.com/aliftoriq/go-crud/initializer"
	"github.com/aliftoriq/go-crud/models"
	"github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
	var body struct {
		Email   string `gorm:"unique"`
		Tittle  string
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
		Tittle:  body.Tittle,
		Content: body.Content,
	}

	result := initializer.DB.Create(&article)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create Article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article Created Succesfuly",
	})
}

func GetArticles(c *gin.Context) {
	var articles []models.Article
	key := "all_article"

	// Attempt to retrieve data from the cache
	art, status, err := cache.GetValueByKey(c, key)
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
	result := initializer.DB.Find(&articles)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get Articles",
			"details": result.Error.Error(),
		})
		return
	}

	// Cache the fetched data
	data, err := json.Marshal(articles)
	if err != nil {
		log.Println("Marshal Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to marshal data for cache",
			"details": err.Error(),
		})
		return
	}

	errSetCache := cache.SetKey(c, key, data, time.Second*3600*24)

	if errSetCache != nil {
		log.Println("Set Cache Error:", errSetCache)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to set cache",
			"details": errSetCache,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    articles,
		"message": "Get Articles Successfully (from database)",
	})
}

func GetArticleByID(c *gin.Context) {
	// Get the article ID from the request parameters
	id := c.Param("id")

	// Define a key to uniquely identify the cached article by its ID
	cacheKey := "article_" + id

	// Attempt to retrieve data from the cache
	art, status, err := cache.GetValueByKey(c, cacheKey)
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

	// Cache miss, fetch data from the database
	var article models.Article
	if result := initializer.DB.First(&article, id); result.Error != nil {
		handleError(c, http.StatusNotFound, "Article not found", result.Error)
		return
	}

	// Cache the fetched data
	data, err := json.Marshal(article)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to marshal article data for cache", err)
		return
	}

	if err := cache.SetKey(c, cacheKey, data, time.Second*60); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to set article cache", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    article,
		"message": "Get Article by ID Successfully (from database)",
	})
}

func UpdateArticle(c *gin.Context) {

	id := c.Param("id")

	var existingArticle models.Article
	if err := initializer.DB.First(&existingArticle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	// Bind the request body to an Article struct
	var updatedArticle models.Article
	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Update the article fields
	existingArticle.Tittle = updatedArticle.Tittle
	existingArticle.Content = updatedArticle.Content

	// Save the updated article to the database
	if err := initializer.DB.Save(&existingArticle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"data":    existingArticle,
	})
}

func DeleteArticle(c *gin.Context) {
	// Get article ID from the request
	id := c.Param("id")

	// Check if the article exists
	var existingArticle models.Article
	if err := initializer.DB.First(&existingArticle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	// Delete the article from the database
	if err := initializer.DB.Delete(&existingArticle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete article",
		})
		return
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
