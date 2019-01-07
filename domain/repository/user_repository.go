package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
)

type UserRepository interface {
	CreateUserByTgUser(user *tgbotapi.User) *models.User
	GetUserByID(id int) *models.User
	RemoveUserByID(id int)
}
