package scraper

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/PuerkitoBio/goquery"
)

var ignoreCategories = []string{"shop", "allergene-zusatzstoffe"}

func ScrapeShop() ([]*Category, error) {
	doc, err := req("shop")
	if err != nil {
		return nil, err
	}

	var categories []*Category

	doc.Find("select.select[name=target]").First().Children().Each(func(i int, s *goquery.Selection) {
		id, ok := s.Attr("value")
		if !ok || slices.Contains(ignoreCategories, id) {
			return
		}

		categories = append(categories, &Category{
			Id:   id,
			Name: s.Text(),
		})
	})

	for _, cat := range categories {
		cat.Items, err = ScrapeCategory(cat.Id)
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func ScrapeCategory(category string) ([]*StoreItem, error) {
	doc, err := req(category)
	if err != nil {
		return nil, err
	}

	var items []*StoreItem
	doc.Find("div.formbody").Each(func(i int, s *goquery.Selection) {
		var item StoreItem

		item.Id, _ = s.Find("input[type=hidden][name=FORM_SUBMIT]").First().Attr("value")
		item.Title = s.Find("h3[itemprop=name]").First().Text()
		item.Description = s.Find("div.description").Text()
		item.Price = s.Find("div.price[itemprop=price]").Text()

		s.Find("fieldset.checkbox_container span").Each(func(i int, selection *goquery.Selection) {
			name, ok := selection.Find("input.checkbox").Attr("name")
			if !ok {
				return
			}
			text := selection.Find("label").Text()
			item.Variants = append(item.Variants, &Variant{Name: name, Description: text})
		})

		dips := s.Find("select[name*=dips] option").Map(func(i int, selection *goquery.Selection) string {
			return selection.Text()
		})
		if len(dips) > 1 {
			item.Dips = dips[1:]
		}

		if !item.IsValid() {
			return
		}

		items = append(items, &item)
	})

	return items, nil
}

func req(path string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://hermans-cafe.de/%s", path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
