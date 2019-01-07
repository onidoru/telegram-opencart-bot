package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"github.com/onidoru/telegram-opencart-bot/opencartsdk"
	"strconv"
	"strings"
	"sync"
)

type Bot struct {
	tgbotapi.BotAPI
}

// NewBot creates and returns new botAPI using the given token.
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &Bot{*bot}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	users := make(map[int]tgbotapi.Message)

	wg := sync.WaitGroup{}
	wg.Add(1)

	for update := range updates {
		userID := getUserID(update)

		// if user is not registered yet, add to map
		if update.Message != nil {
			if _, ok := users[userID]; !ok {
				users[userID] = tgbotapi.Message{}
			}
		} else if update.CallbackQuery != nil {
			if _, ok := users[userID]; !ok {
				users[userID] = tgbotapi.Message{}
			}
		}

		// finally process
		if update.CallbackQuery != nil {
			users[userID], _ = b.processCallback(update, users[userID])
		} else if update.Message != nil {
			users[userID], _ = b.processNewMessage(update)
		}
	}

}

func (b *Bot) processUser() {

}

// processNewMessage processes new user message and returns bot's answer.
func (b Bot) processNewMessage(update tgbotapi.Update) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// process command
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			registerUser(update.Message.From)
			msg.Text = "Hello Sir, Welcome to our Hipster Shop"
		}
	} else {
		msg.Text = "ruh sir no understand but here is our starting keyboard sir"
	}

	// show root category
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	root := c.GetRoot()
	msg.ReplyMarkup = newCategoryKeyboard(root)

	return b.Send(msg)
}

func getUserID(update tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return 0
}

// processCallback processes menu actions and updates existing message.
// No new messages are sent.
func (b Bot) processCallback(update tgbotapi.Update, lastMessage tgbotapi.Message) (tgbotapi.Message, error) {
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	root := c.GetRoot()

	// parse callbackData
	callbackData := update.CallbackQuery.Data

	// assume taxonomy id
	if id, err := strconv.ParseInt(callbackData, 10, 64); err == nil {
		// define if category or item is chosen
		item, category, _ := root.GetByID(id)

		// if item is chosen, show item description and order menu
		if item != nil {
			return b.updateWithItem(lastMessage, item)
		}

		// if category is chosen, show list of all items and subcategories of the chosen category
		if category != nil {
			return b.updateWithCategory(lastMessage, category)
		}
	} else {
		// process the command
		if strings.HasPrefix(callbackData, "back_to_") {
			// find out where to go back
			backTo, _ := strconv.ParseInt(strings.Trim(callbackData, "back_to_"), 10, 64)
			_, category, _ := root.GetByID(backTo)

			return b.updateWithCategory(lastMessage, category)
		} else if strings.HasPrefix(callbackData, "add_to_cart_") {
			// if no cart yet, init cart


			// notify on added item
			alert := tgbotapi.NewCallbackWithAlert("alerted", "Added!")
			alert.CallbackQueryID = update.CallbackQuery.ID
			b.AnswerCallbackQuery(alert)

			return lastMessage, nil
		}
	}

	return tgbotapi.Message{}, nil
}

func registerUser(user *tgbotapi.User) *models.User {
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	c.RegisterUser(user)

	return models.NewUser(user)
}

func updateCart(user models.User) {
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	c.UpdateUserCart(user)
}

func (b Bot) updateWithItem(lastMessage tgbotapi.Message, item *models.Goods) (tgbotapi.Message, error) {

	editedText := tgbotapi.NewEditMessageText(
		lastMessage.Chat.ID,
		lastMessage.MessageID,
		item.String(),
	)
	editedText.ParseMode = tgbotapi.ModeMarkdown

	editedMarkup := tgbotapi.NewEditMessageReplyMarkup(
		lastMessage.Chat.ID,
		lastMessage.MessageID,
		newOrderKeyboard(item.ID, item.ParentID),
	)
	b.Send(editedText)

	return b.Send(editedMarkup)
}

func (b Bot) updateWithCategory(lastMessage tgbotapi.Message, category *models.Category) (tgbotapi.Message, error) {
	editedMarkup := tgbotapi.NewEditMessageReplyMarkup(
		lastMessage.Chat.ID,
		lastMessage.MessageID,
		newCategoryKeyboard(category),
	)

	editedText := tgbotapi.NewEditMessageText(
		lastMessage.Chat.ID,
		lastMessage.MessageID,
		"Menu: ",
	)

	b.Send(editedText)

	return b.Send(editedMarkup)
}
