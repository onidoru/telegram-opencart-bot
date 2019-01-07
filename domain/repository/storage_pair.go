package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
)

type storagePair struct {
	user models.User
	lastMessage tgbotapi.Message
}
