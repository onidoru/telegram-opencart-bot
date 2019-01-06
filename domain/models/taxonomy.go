package models

import (
	"fmt"
	"net/url"
)

type TaxonomyUnit struct {
	ID int64
	ParentID int64

	Name        string
	Description string
	Image       *url.URL
}

func (t *TaxonomyUnit) String() string {
	return fmt.Sprintf("ID: %v, parentID: %v, name: %v, description: %v, image: %v",
		t.ID,
		t.ParentID,
		t.Name,
		t.Description,
		t.Image.String(),
	)
}
