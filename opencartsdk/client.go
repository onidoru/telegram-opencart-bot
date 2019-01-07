package opencartsdk

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/onidoru/telegram-opencart-bot/domain/models"
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

func (client *Client) RegisterUser(user *tgbotapi.User) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		Put("customer")

	if err != nil {
		client.Log.Fatal(err)
	}
	fmt.Println(resp.Body())
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
	user.Cart = parseCart(rawJson.Get("buyItems"))
}

