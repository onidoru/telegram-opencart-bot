package models

import (
	"fmt"
	"github.com/pkg/errors"
)

type Category struct {
	TaxonomyUnit

	Goods           []*Goods
	ChildCategories []*Category
}

func (c Category) String() string {
	return fmt.Sprintf("Category: %v", c.TaxonomyUnit.String())
}

// GetByID returns either category or concrete item that has the given ID.
// Returns an error if not found in the category.
func (c Category) GetByID(id int64) (*Goods, *Category, error) {
	// check if the root is the needed category itself
	if c.ID == id {
		return nil, &c, nil
	}

	// look for concrete items first
	for _, item := range c.Goods {
		if item.ID == id {
			return item, nil, nil
		}
	}

	// if no items found, look for categories
	for _, childCategory := range c.ChildCategories {
		if id == childCategory.ID {
			return nil, childCategory, nil
		}
	}

	// if no categories found directly, traverse recursively
	for _, childCategory := range c.ChildCategories {
		item, foundCategory, _ := childCategory.GetByID(id)
		if (item != nil) || (foundCategory != nil) {
			return item, foundCategory, nil
		}
	}

	// if reached here, nothing is found
	return nil, nil, errors.New("no item with such id is found")
}
