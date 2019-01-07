package models

import (
	"fmt"
)

type Cart struct {
	list map[Goods]int
}

func newCart() *Cart {
	return &Cart{list: make(map[Goods]int)}
}

func (c Cart) AddGoods(item *Goods, amount int) {
	c.list[*item] += amount
}

func (c Cart) RemoveGoods(item *Goods, amount int) {
	c.list [*item] -= amount
}

func (c Cart) String() string {
	var s string
	for item, amount := range c.list {
		s += fmt.Sprintf("☕ %v: %v\n", item, amount)
	}

	return s
}
