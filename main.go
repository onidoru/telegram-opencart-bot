package main

import (
	. "github.com/onidoru/telegram-opencart-bot/bot"
)

func main() {
	bot, err := NewBot()
	if err != nil {
		panic(err)
	}

	bot.Run()
}
