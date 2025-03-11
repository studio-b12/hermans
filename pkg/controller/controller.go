package controller

import (
	"github.com/zekrotja/hermans/pkg/db"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller struct {
	scrapeDb *db.LocalDb[*scraper.Data]
}

func New() (*Controller, error) {
	scrapeDb, err := db.OpenLocalDb[*scraper.Data]("db/scrape_data.msgpack")
	if err != nil {
		return nil, err
	}

	t := &Controller{
		scrapeDb: scrapeDb,
	}
	return t, nil
}

func (t *Controller) Scrape() (*scraper.Data, error) {
	data, err := scraper.ScrapeAll()
	if err != nil {
		return nil, err
	}

	err = t.scrapeDb.Store(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *Controller) GetScrapedData() (*scraper.Data, error) {
	data, err := t.scrapeDb.Load()
	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	return t.Scrape()
}
