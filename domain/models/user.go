package models

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type User struct {
	*tgbotapi.User
	Cart *Cart
}

func NewUser(tgUser *tgbotapi.User) *User  {
	u := &User{User: tgUser}
	u.InitCart()

	return u
}

func (u *User) InitCart()  {
	u.Cart = newCart()
}
