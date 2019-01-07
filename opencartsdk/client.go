package opencartsdk

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
	"github.com/rhymond/go-money"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"net/url"
	"strconv"
)

type Client struct {
	*resty.Client
}

func NewClient(hostURL string) *Client {
	cl := &Client{
		resty.SetHostURL(hostURL),
	}
	cl.SetHostURL(hostURL)

	return cl
}

func (client *Client) GetRoot() *models.Category {
	resp, err := client.R().Get("category/root")
	if err != nil {
		client.Log.Fatal(err)
	}

	rawJson := gjson.ParseBytes(resp.Body())

	rootCategory := &models.Category{}
	rootCategory.ID = rawJson.Get("id").Int()
	rootCategory.Name = rawJson.Get("name").String()
	rootCategory.Description = rawJson.Get("description").String()
	rootCategory.Image, _ = url.Parse(rawJson.Get("image").String())

	parseTaxonomyUnits(rootCategory, rawJson.Get("taxonomyUnits"))

	return rootCategory
}

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

func (client *Client) RegisterUser(user *tgbotapi.User) {
	// set user to form
	formMap := make(map[string]string)
	formMap["is_bot"] = strconv.FormatBool(user.IsBot)
	formMap["first_name"] = user.FirstName
	formMap["last_name"] = user.LastName
	formMap["id"] = strconv.Itoa(user.ID)
	formMap["language_code"] = user.LanguageCode

	// resp, err := client.R().SetFormData(formMap).Put("customer")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		Put("customer")

	if err != nil {
		client.Log.Fatal(err)
	}
	fmt.Println(resp.Body())
}

func (client *Client) GetGoodsByID(id int64) {

}

func (client *Client) UpdateUserCartWithOn(user *models.User, item *models.Goods, amount int) {
	formMap := make(map[string]string)
	formMap["amount"] = strconv.Itoa(amount)

	resp, err := client.R().
		SetFormData(formMap).
		Post("customer/" + strconv.Itoa(user.ID) + "/cart/goods/" + strconv.FormatInt(item.ID, 10))

	if err != nil {
		panic(err)
	}

	client.Log.Println(resp.Body())
}

func (client *Client) UpdateUserCartFromServer(user *models.User) {
	resp, err := client.R().
		Get("customer/" + strconv.Itoa(user.ID) + "/cart/")
	if err != nil {
		client.Log.Fatal(err)
	}

	rawJson := gjson.ParseBytes(resp.Body())
	user.Cart = parseCar(rawJson.Get("buyItems"))
}

func parseCar(rawJson gjson.Result) *models.Cart {
	cart := models.NewCart()
	if rawJson.Get("empty").Bool() == true {
		return cart
	}
	buyItems := rawJson.Get("..0")

	parse := func(key, value gjson.Result) bool {
		item := &models.Goods{}

		if value.Get("goods").Exists() {
			item.ID = value.Get("goos.id").Int()
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
