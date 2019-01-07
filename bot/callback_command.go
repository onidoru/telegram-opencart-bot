package bot

type callbackCommand string

const (
	// remain iddle
	iddle callbackCommand = "iddle"

	// goods menu
	goods_back_to       callbackCommand = "goods_back_to_"
	goods_add_to_cart   callbackCommand = "goods_add_to_cart_"
	goods_item_info     callbackCommand = "goods_show_item_"
	goods_next_category callbackCommand = "goods_next_category"

	// main menu
	to_main_menu callbackCommand = "main_to_main"
	to_root      callbackCommand = "main_to_root"
	to_cart_menu callbackCommand = "main_to_cart_menu"
	to_settings  callbackCommand = "main_to_settings"

	// cart menu
	cart_view      callbackCommand = "cart_view"
	cart_purchase  callbackCommand = "cart_purchase"
	cart_drop      callbackCommand = "cart_drop"
	cart_incr_item callbackCommand = "cart_incr_item_"
	cart_decr_item callbackCommand = "cart_decr_item_"
)

func (c callbackCommand) String() string {
	return string(c)
}
