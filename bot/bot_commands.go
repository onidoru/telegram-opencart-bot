package bot

import (
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"github.com/rhymond/go-money"
	url2 "net/url"
	"strconv"
	"strings"
)

type botCommand string

const (
	start    botCommand = "start"
	add_item botCommand = "additem"
)

type botCommandArg string

const (
	pref botCommandArg = "-"
	suf  botCommandArg = "="

	id          botCommandArg = "-id"
	name        botCommandArg = "-name="
	description botCommandArg = "-description="
	image       botCommandArg = "-image="
	price       botCommandArg = "-price="
)

func (b botCommand) String() string {
	return string(b)
}

func (b botCommandArg) String() string {
	return string(b)
}

func parseNewItemCommand(s string) *models.Goods {
	item := &models.Goods{}

	unparsedArgs := strings.Fields(s)

	for _, unparsedArg := range unparsedArgs {
		if strings.Contains(unparsedArg, image.String()) {
			trim := strings.Trim(unparsedArg, image.String())
			item.Image, _ = url2.Parse(trim)
			continue
		}
		if strings.Contains(unparsedArg, name.String()) {
			trim := strings.Trim(unparsedArg, name.String())
			item.Name = trim
			continue
		}

		if strings.Contains(unparsedArg, price.String()) {
			trim := strings.Trim(unparsedArg, price.String())
			amount, _ := strconv.ParseInt(trim, 10, 64)
			item.Price = *money.New(amount, "UAH")
			continue
		}

		if strings.Contains(unparsedArg, description.String()) {
			trim := strings.Trim(unparsedArg, description.String())
			item.Description = trim
			continue
		}
	}

	return item
}
