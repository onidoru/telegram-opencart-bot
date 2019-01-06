package main

import (
	. "github.com/onidoru/telegram-opencart-bot/bot"
)

func main() {
	// Initialize bot with the given token.
	bot, err := NewBot("332637329:AAGjEdLxDveCbukkhz-7htYngej5vTrjTws")
	if err != nil {
		panic(err)
	}

	bot.Run()
}
