package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/zekrotja/hermans/pkg/scraper"
)

func main() {
	// categories, err := scraper.ScrapeShop()
	// if err != nil {
	// 	panic(err)
	// }

	// spew.Dump(categories)

	items, err := scraper.ScrapeCategory("currywurst-snacks")
	if err != nil {
		panic(err)
	}

	spew.Dump(items)
}
