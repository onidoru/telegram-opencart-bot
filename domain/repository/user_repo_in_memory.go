package repository

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
)

type UserRepositoryInMemory struct {
	m map[int]*models.User
}

func (c *UserRepositoryInMemory) CreateUserByTgUser(user *tgbotapi.User) *models.User {
	u := models.NewUser(user)
	c.m[u.ID] = u

	return u
}

func (c *UserRepositoryInMemory) GetUserByID(id int) *models.User {
	return c.m[id]
}

func (c *UserRepositoryInMemory) RemoveUserByID(id int) {
	delete(c.m, id)
}
