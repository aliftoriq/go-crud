package repositories

import (
	"github.com/aliftoriq/go-crud/initializer"
	"github.com/aliftoriq/go-crud/models"
	"gorm.io/gorm"
)

//go:generate mockery --outpkg mocks --name UserRepository
type UserRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: initializer.DB}
}

func (ur *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepository) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	err := ur.db.First(user, id).Error
	return user, err
}

func (ur *userRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := ur.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (ur *userRepository) Update(user *models.User) error {
	return ur.db.Save(user).Error
}

func (ur *userRepository) Delete(user *models.User) error {
	return ur.db.Unscoped().Delete(user).Error
}
