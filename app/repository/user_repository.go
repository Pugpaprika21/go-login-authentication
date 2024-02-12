package repository

import (
	"strings"
	"time"

	"github.com/Pugpaprika21/go-fiber/app/db"
	"github.com/Pugpaprika21/go-fiber/app/models"
	"github.com/Pugpaprika21/go-fiber/dto"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Save(body dto.UserRegisterBodyRequest) (bool, error)
	GetByID(id string) (*dto.UserQueryRow, error)
	Get() ([]dto.UserQueryRow, error)
	UsernameIsExiting(username string) int64
	CheckLogin(username string, password string) (*dto.UserQueryRow, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *userRepository {
	return &userRepository{
		DB: db.Conn,
	}
}

func (u *userRepository) Save(body dto.UserRegisterBodyRequest) (bool, error) {
	result := u.DB.Model(&models.User{}).Create(map[string]any{
		"CreatedAt": time.Now(),
		"Username":  body.Username,
		"Password":  body.Password,
		"Email":     body.Email,
		"Token":     body.Token,
	})

	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (u *userRepository) GetByID(id string) (*dto.UserQueryRow, error) {
	var row dto.UserQueryRow
	result := u.DB.Model(&models.User{}).Where("id = ?", id).Find(&row)
	return &row, result.Error
}

func (u *userRepository) Get() ([]dto.UserQueryRow, error) {
	var rows []dto.UserQueryRow
	result := u.DB.Model(&models.User{}).Find(&rows)
	return rows, result.Error
}

func (u *userRepository) UsernameIsExiting(username string) int64 {
	var num int64
	u.DB.Model(&models.User{}).Where("username = ?", strings.TrimSpace(username)).Count(&num)
	return num
}

func (u *userRepository) CheckLogin(username, password string) (*dto.UserQueryRow, error) {
	var row dto.UserQueryRow

	result := u.DB.Model(&models.User{}).Where("username = ?", strings.TrimSpace(username)).First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return &row, nil
}
