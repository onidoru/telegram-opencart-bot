package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
)

type userInMemoryStorage struct {
	m map[int]storagePair
}

func (c *userInMemoryStorage) GetLastMessageByID(id int) (tgbotapi.Message, bool) {
	pair, ok := c.m[id]
	return pair.lastMessage, ok
}

func (c *userInMemoryStorage) StoreLastMessageByID(id int, message tgbotapi.Message) {
	pair := c.m[id]
	pair.lastMessage = message
	c.m[id] = pair
}

func newInMemoryStorage() *userInMemoryStorage {
	return &userInMemoryStorage{m: make(map[int]storagePair)}
}

func (c *userInMemoryStorage) CreateWithTgUser(user *tgbotapi.User) *models.User {
	s := storagePair{
		user: models.NewUser(user),
	}
	c.m[s.user.ID] = s

	return s.user
}

func (c *userInMemoryStorage) GetUserByID(id int) (*models.User, bool) {
	sPair, ok := c.m[id]
	return sPair.user, ok
}

func (c *userInMemoryStorage) RemoveUserByID(id int) {
	delete(c.m, id)
}

