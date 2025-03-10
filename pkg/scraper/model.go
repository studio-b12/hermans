package scraper

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
}
