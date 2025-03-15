package controller

import (
	"time"

	"github.com/google/uuid"
	"github.com/zekrotja/hermans/pkg/cache"
	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller struct {
	scrapeCache *cache.LocalCache[*scraper.Data]
}

func New() (*Controller, error) {
	scrapeDb, err := cache.OpenLocalCache[*scraper.Data]("db/scrape_data.msgpack")
	if err != nil {
		return nil, err
	}

	t := &Controller{
		scrapeCache: scrapeDb,
	}
	return t, nil
}

func (t *Controller) Scrape() (*scraper.Data, error) {
	data, err := scraper.ScrapeAll()
	if err != nil {
		return nil, err
	}

	err = t.scrapeCache.Store(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *Controller) GetScrapedData() (*scraper.Data, error) {
	data, err := t.scrapeCache.Load()
	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	return t.Scrape()
}

func (t *Controller) CreateOrderList() {

}
