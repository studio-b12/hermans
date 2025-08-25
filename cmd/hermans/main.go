package main

import (
	"log/slog"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	"github.com/zekrotja/hermans/pkg/api"
	"github.com/zekrotja/hermans/pkg/controller"
	"github.com/zekrotja/hermans/pkg/database"
)

type Args struct {
	BindAddress    string     `arg:"--bind-address,env:HMS_BIND_ADDRESS" help:"Address to bind to" default:"0.0.0.0:8080"`
	DatabaseDsn    string     `arg:"--database-dsn,required,env:HMS_DATABASE_DSN" help:"Database DSN"`
	CacheDir       string     `arg:"--cache-dir,env:HMS_CACHE_DIR" help:"Cache directory" default:"./cache"`
	LogLevel       slog.Level `arg:"--log-level,env:HMS_LOG_LEVEL" help:"Log level" default:"info"`
	ScrapeInterval string     `arg:"--scrape-interval,env:HMS_SCRAPE_INTERVAL" help:"Interval for periodic scraping (e.g., '1h', '30m')" default:"168h"`
}

func checkErr(msg string, err error, extraFields ...any) {
	if err == nil {
		return
	}

	fields := append([]any{"err", err}, extraFields...)
	slog.Error(msg, fields...)
	os.Exit(1)
}

func main() {
	godotenv.Load()

	var args Args
	arg.MustParse(&args)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: args.LogLevel}))
	slog.SetDefault(logger)

	slog.Info("initializing database connection ...", "dsn", args.DatabaseDsn)
	db, err := database.New(args.DatabaseDsn)
	checkErr("failed initializing database", err)

	slog.Info("initializing controller ...")
	ctl, err := controller.New(args.CacheDir, db)
	checkErr("failed initializing controller", err)

	slog.Info("starting scraping scheduler ...", "interval", args.ScrapeInterval)
	go ctl.StartScrapingScheduler(args.ScrapeInterval)

	a := api.New(ctl, args.BindAddress)

	slog.Info("starting web server ...", "addr", args.BindAddress)
	err = a.Start()
	checkErr("failed starting web server", err)
}
