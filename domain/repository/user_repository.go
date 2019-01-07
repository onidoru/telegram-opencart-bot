package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
)

type UserRepository interface {
	CreateWithTgUser(user *tgbotapi.User) *models.User
	GetUserByID(id int) (*models.User, bool)
	RemoveUserByID(id int)
	GetLastMessageByID(id int) (tgbotapi.Message, bool)
	StoreLastMessageByID(id int, message tgbotapi.Message)
}

func NewInMemoryStorage() UserRepository {
	return UserRepository(newInMemoryStorage())
}

