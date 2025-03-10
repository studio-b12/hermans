package scraper

type Variant struct {
	Name        string
	Description string
}

type Category struct {
	Id    string
	Name  string
	Items []*StoreItem
}

type StoreItem struct {
	Id          string
	Title       string
	Description string
	Price       string
	Variants    []*Variant
	Dips        []string
}

func (t *StoreItem) IsValid() bool {
	return t.Id != "" && t.Title != ""
}
