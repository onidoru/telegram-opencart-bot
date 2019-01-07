package bot

import (
	"strconv"
	"strings"
)

type callbackCommand string

const (
	// remain iddle
	iddle callbackCommand = "iddle"

	// goods menu
	goods_prefix        callbackCommand = "goods_"
	goods_back_to       callbackCommand = "goods_back_to_"
	goods_add_to_cart   callbackCommand = "goods_add_to_cart_"
	goods_item_info     callbackCommand = "goods_show_item_"
	goods_next_category callbackCommand = "goods_next_category_"

	// main menu
	main_to_prefix    callbackCommand = "main_to"
	main_to_menu      callbackCommand = "main_to_menu"
	main_to_root      callbackCommand = "main_to_root"
	main_to_cart_menu callbackCommand = "main_to_cart_menu"
	main_to_settings  callbackCommand = "main_to_settings"

	// cart menu
	cart_prefix    callbackCommand = "cart_"
	cart_view      callbackCommand = "cart_view"
	cart_purchase  callbackCommand = "cart_purchase"
	cart_drop      callbackCommand = "cart_drop"
	cart_incr_item callbackCommand = "cart_incr_item_"
	cart_decr_item callbackCommand = "cart_decr_item_"
)

func (c callbackCommand) String() string {
	return string(c)
}

// parseCallbackCommand returns command and its argument
func parseCallbackCommand(s string) (callbackCommand, int64) {
	// quick determine if iddle:
	if strings.EqualFold(s, iddle.String()) {
		return iddle, 0
	}

	// find prefix

	// determine if main menu
	switch s {
	case main_to_menu.String():
		return main_to_menu, 0
	case main_to_root.String():
		return main_to_root, 0
	case main_to_cart_menu.String():
		return main_to_cart_menu, 0
	case main_to_settings.String():
		return main_to_settings, 0
	case cart_view.String():
		return cart_view, 0
	case cart_purchase.String():
		return cart_purchase, 0
	case cart_drop.String():
		return cart_drop, 0
	}

	// determine if cart menu
	if strings.Contains(s, cart_prefix.String()) {
		// parse commands with suffix
		if strings.Contains(s, cart_incr_item.String()) {
			t := strings.Trim(s, cart_incr_item.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return cart_incr_item, id
		}
		if strings.Contains(s, cart_decr_item.String()) {
			t := strings.Trim(s, cart_decr_item.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return cart_decr_item, id
		}
	}

	// determine if goods menu
	if strings.Contains(s, goods_prefix.String()) {

		// parse commands with suffix
		if strings.Contains(s, goods_back_to.String()) {
			t := strings.Trim(s, goods_back_to.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return goods_back_to, id
		}
		if strings.Contains(s, goods_item_info.String()) {
			t := strings.Trim(s, goods_item_info.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return goods_item_info, id
		}
		if strings.Contains(s, goods_next_category.String()) {
			t := strings.Trim(s, goods_next_category.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return goods_next_category, id
		}
		if strings.Contains(s, goods_add_to_cart.String()) {
			t := strings.Trim(s, goods_add_to_cart.String())
			id, _ := strconv.ParseInt(t, 10, 64)
			return goods_add_to_cart, id
		}
	}

	return main_to_menu, 0
}
