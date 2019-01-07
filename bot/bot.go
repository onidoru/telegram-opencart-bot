package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"github.com/onidoru/telegram-opencart-bot/domain/repository"
	"github.com/onidoru/telegram-opencart-bot/opencartsdk"
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

	users := repository.NewInMemoryStorage()

	for update := range updates {
		userID := getUserID(update)

		// if user is not registered yet, add to map
		if update.Message != nil {
			if _, ok := users.GetUserByID(userID); !ok {
				users.CreateWithTgUser(update.Message.From)
				users.StoreLastMessageByID(userID, tgbotapi.Message{})
			}
		} else if update.CallbackQuery != nil {
			if _, ok := users.GetUserByID(userID); !ok {
				users.CreateWithTgUser(update.CallbackQuery.From)
				users.StoreLastMessageByID(userID, tgbotapi.Message{})
			}
		} else if update.PreCheckoutQuery != nil {
			newPrecheckoutConfig := tgbotapi.PreCheckoutConfig{
				PreCheckoutQueryID: update.PreCheckoutQuery.ID,
				OK:                 true,
				ErrorMessage:       "",
			}

			responce, err := b.AnswerPreCheckoutQuery(newPrecheckoutConfig)

			fmt.Println(responce, err)
		}

		// finally process
		if update.CallbackQuery != nil {
			lastMessage, _ := users.GetLastMessageByID(userID)
			user, _ := users.GetUserByID(userID)
			lastMessage, _ = b.processCallback(update, lastMessage, user)
			users.StoreLastMessageByID(userID, lastMessage)
		} else if update.Message != nil {
			lastMessage, _ := b.processNewMessage(update)
			users.StoreLastMessageByID(userID, lastMessage)
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
			msg.Text = "Hello Sir, Welcome to our Hipster Shop!"
		}
	} else {
		registerUser(update.Message.From)
		msg.Text = "Please Sir, user our Menu!"
	}

	msg.ReplyMarkup = newMainMenuKeyboard()

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
func (b Bot) processCallback(update tgbotapi.Update, lastMessage tgbotapi.Message, user *models.User) (tgbotapi.Message, error) {

	// parse callback command
	fmt.Println(update.CallbackQuery.Data)
	cbckCommand, arg := parseCallbackCommand(update.CallbackQuery.Data)
	// callbackData := update.CallbackQuery.Data

	// update all items, initialize root
	client := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	root := client.GetRoot()

	var item *models.Goods
	var category *models.Category

	// find item or category on the given arg
	if arg != 0 {
		var err error
		item, category, err = root.GetByID(arg)
		if err != nil {
			cbckCommand = iddle
		}
	}

	switch cbckCommand {
	case iddle:
	// stay iddle, do nothing

	case goods_item_info:
		// show item info
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				item.String(),
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newOrderKeyboard(item.ID, item.ParentID),
			),
		)

	case goods_back_to:
		// back to the non-root category
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"Menu: ",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCategoryKeyboard(category),
			),
		)

	case goods_add_to_cart:
		// add chosen item to cart
		// get chosen item id and concrete item from root
		if user.Cart == nil {
			user.InitCart()
		}

		user.Cart.AddGoods(item, 1)
		updateCart(user, item, 1)

		// notify on added item
		alert := tgbotapi.NewCallbackWithAlert("alerted", "Added!")
		alert.CallbackQueryID = update.CallbackQuery.ID
		b.AnswerCallbackQuery(alert)

		return lastMessage, nil

	case goods_next_category:
		// process to the chosen child category
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"Menu: ",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCategoryKeyboard(category),
			),
		)

	case main_to_menu:
		// return to main menu
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"Hello Sir, Welcome to our Hipster Shop",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newMainMenuKeyboard(),
			),
		)

	case main_to_root:
		// process to the root category and show list of goods
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"Menu: ",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCategoryKeyboard(root),
			),
		)

	case main_to_cart_menu:
		// show cart menu
		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"🛒 Cart:",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCartMenuKeyboard(),
			),
		)

	case main_to_settings:
	// stay iddle for now

	case cart_view:
		// show the cart content
		client.UpdateUserCartFromServer(user)

		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"🛒 Cart:",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCartViewKeyboard(user.Cart),
			),
		)

	case cart_purchase:
	// TODO: Cart Purchase
	// stay iddle for now

	case cart_drop:
	// TODO: Cart Drop
	// stay iddle for now

	case cart_incr_item:
		// update cart and load new cart
		user.Cart.AddGoods(item, 1)                // add locally
		client.UpdateUserCartWithOn(user, item, 1) // add remotely

		// update from server to make sure it's good
		client.UpdateUserCartFromServer(user)

		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"🛒 Cart:",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCartViewKeyboard(user.Cart),
			),
		)

	case cart_decr_item:
		// update cart and load new cart

		user.Cart.RemoveGoods(item, 1)        // remove locally
		client.RemoveItemFromCart(user, item) // remove remotely

		// update from server to make sure it's good
		client.UpdateUserCartFromServer(user)

		return b.updateWith(
			tgbotapi.NewEditMessageText(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				"🛒 Cart:",
			),
			tgbotapi.NewEditMessageReplyMarkup(
				lastMessage.Chat.ID,
				lastMessage.MessageID,
				newCartViewKeyboard(user.Cart),
			),
		)
	}
	//
	// // assume taxonomy id
	// if id, err := strconv.ParseInt(callbackData, 10, 64); err == nil {
	// 	// define if category or item is chosen
	// 	item, category, _ := root.GetByID(id)
	//
	// 	// if item is chosen, show item description and order menu
	// 	if item != nil {
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				item.String(),
	// 			),
	//
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newOrderKeyboard(item.ID, item.ParentID),
	// 			),
	// 		)
	// 	}
	//
	// 	// if category is chosen, show list of all items and subcategories of the chosen category
	// 	if category != nil {
	//
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				"Menu: ",
	// 			),
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newCategoryKeyboard(category),
	// 			),
	// 		)
	// 	}
	// } else {
	// 	// process commands
	//
	// 	// found main_to_root command
	// 	if strings.HasPrefix(callbackData, "main_to_root") {
	// 		// show root menu
	// 		c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	// 		root := c.GetRoot()
	//
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				"Menu: ",
	// 			),
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newCategoryKeyboard(root),
	// 			),
	// 		)
	// 	}
	//
	// 	// found main_to_cart_menu command
	// 	if strings.EqualFold(callbackData, "main_to_cart_menu") {
	// 		// update user view with cart menu
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				"🛒 Cart:",
	// 			),
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newCartMenuKeyboard(),
	// 			),
	// 		)
	//
	// 	}
	//
	// 	// found cart_view command
	// 	if strings.EqualFold(callbackData, "cart_view") {
	// 		// show cart menu
	// 		c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	// 		c.UpdateUserCartFromServer(user)
	//
	// 		// fmt.Println(user.Cart)
	//
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				"🛒 Cart:",
	// 			),
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newCartViewKeyboard(user.Cart),
	// 			),
	// 		)
	// 	}
	//
	// 	// found goods_back_to command from concrete item description
	// 	if strings.HasPrefix(callbackData, "back_to_") {
	// 		// find out where to go back
	// 		cutted := strings.Trim(callbackData, "back_to_")
	// 		backTo, err := strconv.ParseInt(cutted, 10, 64)
	// 		fmt.Println(cutted)
	// 		if err != nil {
	// 			switch cutted {
	// 			case "menu":
	// 				return b.updateWith(
	// 					tgbotapi.NewEditMessageText(
	// 						lastMessage.Chat.ID,
	// 						lastMessage.MessageID,
	// 						"Hello Sir, Welcome to our Hipster Shop",
	// 					),
	// 					tgbotapi.NewEditMessageReplyMarkup(
	// 						lastMessage.Chat.ID,
	// 						lastMessage.MessageID,
	// 						newMainMenuKeyboard(),
	// 					),
	// 				)
	// 			case "kek":
	// 			}
	// 		}
	// 		_, category, _ := root.GetByID(backTo)
	//
	// 		return b.updateWith(
	// 			tgbotapi.NewEditMessageText(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				"Menu: ",
	// 			),
	// 			tgbotapi.NewEditMessageReplyMarkup(
	// 				lastMessage.Chat.ID,
	// 				lastMessage.MessageID,
	// 				newCategoryKeyboard(category),
	// 			),
	// 		)
	// 	} else if strings.HasPrefix(callbackData, "add_to_cart_") {
	// 		// get chosen item id and concrete item from root
	// 		itemID, _ := strconv.ParseInt(strings.Trim(callbackData, "add_to_cart_"), 10, 64)
	// 		item, _, _ := root.GetByID(itemID)
	//
	// 		if user.Cart == nil {
	// 			user.InitCart()
	// 		}
	//
	// 		user.Cart.AddGoods(item, 1)
	// 		updateCart(user, item, 1)
	//
	// 		// notify on added item
	// 		alert := tgbotapi.NewCallbackWithAlert("alerted", "Added!")
	// 		alert.CallbackQueryID = update.CallbackQuery.ID
	// 		b.AnswerCallbackQuery(alert)
	//
	// 		return lastMessage, nil
	// 	}
	// }

	return tgbotapi.Message{}, nil
}

func registerUser(user *tgbotapi.User) *models.User {
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	c.RegisterUser(user)

	return models.NewUser(user)
}

func updateCart(user *models.User, item *models.Goods, amount int) {
	c := opencartsdk.NewClient("https://telegram-coffee-shop.herokuapp.com/")
	c.UpdateUserCartWithOn(user, item, amount)
}

func (b Bot) updateWith(
	editedText tgbotapi.EditMessageTextConfig,
	editedMarkup tgbotapi.EditMessageReplyMarkupConfig,
) (tgbotapi.Message, error) {
	editedText.ParseMode = tgbotapi.ModeMarkdown
	b.Send(editedText)

	return b.Send(editedMarkup)
}
