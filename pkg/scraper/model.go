package scraper

type Data struct {
	Categories []*Category `json:"categories"`
	Drinks     []*Drink    `json:"drinks"`
}

type Variant struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Category struct {
	Id    string       `json:"id"`
	Name  string       `json:"name"`
	Items []*StoreItem `json:"items"`
}

type StoreItem struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       string     `json:"price"`
	Variants    []*Variant `json:"variants"`
	Dips        []string   `json:"dips"`
}

func (t *StoreItem) VariantsContain(name string) bool {
	for _, variant := range t.Variants {
		if variant.Name == name {
			return true
		}
	}
	return false
}

func (t *StoreItem) IsValid() bool {
	return t.Id != "" && t.Title != ""
}

type Drink struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}
