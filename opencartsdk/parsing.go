package opencartsdk

import (
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"github.com/rhymond/go-money"
	"github.com/tidwall/gjson"
	"net/url"
)

func parseTaxonomyUnits(rootCategory *models.Category, jsonUnits gjson.Result) {

	parse := func(key, value gjson.Result) bool {

		if value.Get("price").Exists() {
			g := &models.Goods{}
			g.ID = value.Get("id").Int()
			g.ParentID = rootCategory.ID
			g.Name = value.Get("name").String()
			g.Description = value.Get("description").String()
			g.Price = *money.New(value.Get("price").Int(), "UAH")
			g.Image, _ = url.Parse(value.Get("image").String())
			rootCategory.Goods = append(rootCategory.Goods, g)
		} else if value.Get("taxonomyUnits").Exists() {
			c := &models.Category{}

			c.ID = value.Get("id").Int()
			c.ParentID = rootCategory.ID
			c.Name = value.Get("name").String()
			c.Description = value.Get("description").String()
			c.Image, _ = url.Parse(value.Get("image").String())

			parseTaxonomyUnits(c, value.Get("taxonomyUnits"))

			rootCategory.ChildCategories = append(rootCategory.ChildCategories, c)

		}

		return true
	}

	jsonUnits.ForEach(parse)
}

func parseCart(rawJson gjson.Result) *models.Cart {
	cart := models.NewCart()
	if rawJson.Get("empty").Bool() == true {
		return cart
	}
	buyItems := rawJson.Get("..0")

	parse := func(key, value gjson.Result) bool {
		item := &models.Goods{}

		if value.Get("goods").Exists() {
			item.ID = value.Get("goods.id").Int()
			item.Name = value.Get("goods.name").String()
			item.Description = value.Get("goods.description").String()
			item.Image, _ = url.Parse(value.Get("goods.image").String())
			item.Price = *money.New(value.Get("goods.price").Int(), "UAH")
		}

		amount := int(value.Get("amount").Int())
		cart.AddGoods(item, amount)

		return true
	}

	buyItems.ForEach(parse)

	return cart
}
