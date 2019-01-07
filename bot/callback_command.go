package bot

type callbackCommand string

const (
	// remain iddle
	iddle callbackCommand = "iddle"

	// goods menu
	back_to     callbackCommand = "goods_back_to_"
	add_to_cart callbackCommand = "goods_add_to_cart_"
	item_info   callbackCommand = "goods_show_item_"

	// main menu
	to_main_menu callbackCommand = "main_to_main"
	to_root      callbackCommand = "main_to_root"
	to_cart_menu callbackCommand = "main_to_cart_menu"
	to_settings  callbackCommand = "main_to_settings"

	// cart menu
	view_cart     callbackCommand = "cart_view"
	purchase_cart callbackCommand = "cart_purchase"
	drop_cart     callbackCommand = "cart_drop"
	incr_item     callbackCommand = "cart_incr_item_"
	decr_item     callbackCommand = "cart_decr_item_"
)

func (c callbackCommand) String() string {
	return string(c)
}
