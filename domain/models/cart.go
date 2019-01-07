package models

import (
	"fmt"
)

type Cart struct {
	isEmpty bool
	list    map[Goods]int
}

func NewCart() *Cart {
	return &Cart{list: make(map[Goods]int)}
}

func (c *Cart) IsEmpty() bool {
	return c.isEmpty
}

// true is empty, false is not
func (c *Cart) updateEmpty() {
	c.isEmpty = len(c.list) == 0
}

func (c *Cart) AddGoods(item *Goods, amount int) {
	c.list[*item] += amount
	c.updateEmpty()
}

func (c *Cart) RemoveGoods(item *Goods, amount int) {
	c.list [*item] -= amount
	c.updateEmpty()
}

func (c *Cart) String() string {
	var s string

	for item, amount := range c.list {
		s += fmt.Sprintf("from cart: %v: %v\n", item.String(), amount)
	}

	return s
}

func (c *Cart) GetAllGoods() map[Goods]int {
	return c.list
}

func (c *Cart) CountTotalAmount() int {
	if c.isEmpty {
		return 0
	}

	price := 0
	for item, amount := range c.list {
		price += int(item.Price.Amount()) * amount
	}

	return price
}
