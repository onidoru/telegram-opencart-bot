package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"strconv"
)

func newCategoryKeyboard(rootCategory *models.Category) tgbotapi.InlineKeyboardMarkup {
	// make buttons for childCategory goods
	var categoryItemButtons [][]tgbotapi.InlineKeyboardButton

	for _, item := range rootCategory.Goods {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"☕ "+item.Name,
				goods_item_info.String()+strconv.FormatInt(item.ID, 10),
			),
		)
		categoryItemButtons = append(categoryItemButtons, row)
	}

	// make buttons for children of the given childCategory
	for _, childCategory := range rootCategory.ChildCategories {
		if childCategory.ParentID == rootCategory.ID {
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"📒 "+childCategory.Name,
					goods_next_category.String()+strconv.FormatInt(childCategory.ID, 10)),
			)
			categoryItemButtons = append(categoryItemButtons, row)
		}
	}

	// add back to start menu button if the category is not root
	if rootCategory.ParentID != 0 {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back",
				goods_back_to.String()+strconv.FormatInt(rootCategory.ParentID, 10),
			),
		)

		categoryItemButtons = append(categoryItemButtons, row)
	} else {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back To Main Menu",
				main_to_menu.String(),
			),
		)

		categoryItemButtons = append(categoryItemButtons, row)
	}

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(categoryItemButtons...)

	return replyKeyboard
}

func newOrderKeyboard(itemID, backTo int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Order Item",
				goods_add_to_cart.String()+strconv.FormatInt(itemID, 10),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"Back",
				goods_back_to.String()+strconv.FormatInt(backTo, 10),
			),
		),
	)
}

func newMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"📃 Menu",
				main_to_root.String(),
			)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🛒 Cart",
				main_to_cart_menu.String(),
			)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"⚙️ Settings",
				main_to_settings.String(),
			)),
	)
}

func newCartViewKeyboard(cart *models.Cart) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// loop the cart
	for item, amount := range cart.GetAllGoods() {
		// for each item create two rows:
		// first one with item name + amount + total cost
		// second one with +/- buttons

		fmt.Println(item, amount)

		itemRow := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(
					"☕️ %v x %v - %v",
					item.Name, amount, item.Price.Amount(),
				),
				iddle.String(),
			),
		)

		regAmountRow := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"➕",
				cart_incr_item.String()+strconv.FormatInt(item.ID, 10),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"➖",
				cart_decr_item.String()+strconv.FormatInt(item.ID, 10),
			),
		)

		rows = append(rows, itemRow)
		rows = append(rows, regAmountRow)
	}

	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			"🛒  Back to Cart",
			main_to_cart_menu.String(),
		))

	rows = append(rows, backRow)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return keyboard
}

func newCartMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🛒 View Cart",
				cart_view.String(),
			)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"💰 Purchase",
				cart_purchase.String(),
			)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✖️️ Drop Cart",
				cart_drop.String(),
			)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🔙 Back to Main Menu",
				main_to_menu.String(),
			)),
	)
}
