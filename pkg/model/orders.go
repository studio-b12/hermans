package model

type StoreItem struct {
	Id       string   `json:"id"`
	Variants []string `json:"variants"`
	Dips     []string `json:"dips"`
}

type Drink struct {
	Name string `json:"name"`
}

type Order struct {
	Id        string     `json:"id"`
	StoreItem *StoreItem `json:"store_item"`
	Drink     *Drink     `json:"drink"`
}
