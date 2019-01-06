package models

import (
	"fmt"
	"github.com/rhymond/go-money"
)

type Goods struct {
	TaxonomyUnit
	Price money.Money
}

func (g *Goods) String() string {

	s := fmt.Sprintf("☕ *%v*, ", g.Name)
	s += fmt.Sprintf(" _%v %v_\n", g.Price.Amount(), g.Price.Currency().Code)
	s += fmt.Sprintf("%v\n", g.Description)
	if g.Image.String() != "" {
		s += fmt.Sprintf("[Photo](%v)", g.Image.String())
	}

	return s
}
