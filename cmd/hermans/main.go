package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/zekrotja/hermans/pkg/controller"
)

func main() {
	ctl, err := controller.New()
	if err != nil {
		panic(err)
	}

	data, err := ctl.GetScrapedData()
	if err != nil {
		panic(err)
	}

	spew.Dump(data)
}
