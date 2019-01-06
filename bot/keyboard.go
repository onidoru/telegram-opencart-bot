package bot

import (
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
				strconv.FormatInt(item.ID, 10),
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
					strconv.FormatInt(childCategory.ID, 10)),
			)
			categoryItemButtons = append(categoryItemButtons, row)
		}
	}

	// add back to start menu button if the category is not root
	if rootCategory.ParentID != 0 {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Back",
				"back_to_"+strconv.FormatInt(rootCategory.ParentID, 10),
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
				"order_item_"+strconv.FormatInt(itemID, 10),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"Back",
				"back_to_"+strconv.FormatInt(backTo, 10),
			),
		),
	)
}
