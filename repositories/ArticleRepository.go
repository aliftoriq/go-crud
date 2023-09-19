package repositories

import (
	"errors"

	"github.com/aliftoriq/go-crud/initializer"
	"github.com/aliftoriq/go-crud/models"
	"gorm.io/gorm"
)

//go:generate mockery --outpkg mocks --name ArticleRepository
type ArticleRepository interface {
	CreateArticle(article models.Article) error
	GetArticles() (*[]models.Article, error)
	GetArticleById(id string) (*models.Article, error)
	UpdateArticle(id string, article models.Article) error
	DeleteArticle(id string) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository() ArticleRepository {
	return &articleRepository{db: initializer.DB}
}

func (ar *articleRepository) CreateArticle(article models.Article) error {
	return ar.db.Create(&article).Error
}

func (ar *articleRepository) GetArticles() (*[]models.Article, error) {
	var article []models.Article
	art := ar.db.Find(&article)

	if art.Error != nil {
		return nil, errors.New("FAILED TO GET ARTICLES")
	}

	return &article, nil
}

func (ar *articleRepository) GetArticleById(id string) (*models.Article, error) {
	var article models.Article
	art := ar.db.Find(&article, id)

	if art.Error != nil {
		return nil, errors.New("FAILED TO GET ARTICLES")
	}

	return &article, nil
}

func (ar *articleRepository) UpdateArticle(id string, article models.Article) error {
	var existingArticle models.Article
	if err := initializer.DB.First(&existingArticle, id).Error; err != nil {
		return errors.New("ARTICLE NOT FOUND")
	}

	existingArticle.Title = article.Title
	existingArticle.Content = article.Content

	if err := initializer.DB.Save(&existingArticle).Error; err != nil {
		return errors.New("FAILED TO UPDATE ARTICLE")
	}

	return nil
}

func (ar *articleRepository) DeleteArticle(id string) error {
	var existingArticle models.Article
	if err := ar.db.First(&existingArticle, id).Error; err != nil {
		return errors.New("ARTICLE NOT FOUND")
	}

	if err := ar.db.Delete(&existingArticle).Error; err != nil {
		return errors.New("FAILED TO DELETE ARTICLE")
	}

	return nil
}
