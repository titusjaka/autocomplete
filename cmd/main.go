package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gosexy/to"
	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/titusjaka/autocomplete/database/bolt"
	"github.com/titusjaka/autocomplete/proxy"
	"github.com/titusjaka/autocomplete/server"
)

func main() {
	// nolint:lll
	var opts = struct {
		Listen                 string `long:"http.listen" env:"AVS_STUB_LISTEN" default:":8080" description:"HTTP Listen address"`
		Verbose                bool   `long:"verbose" env:"AVS_STUB_VERBOSE" description:"Verbose output"`
		DBPath                 string `long:"db.path" env:"AVS_STUB_DBPATH" default:"autocomplete.db" description:"Path to bolt DB file"`
		AVSAutocompleteURL     string `long:"avs.autocomplete.url" env:"AVS_AUTOCOMPLETE_URL" default:"https://places.aviasales.ru" description:"URL of Aviasales autocomplete service"`
		AVSAutocompleteTimeout string `long:"avs.autocomplere.timeout" env:"AVS_AUTOCOMPLETE_TIMEOUT" default:"3s" description:"Timeout for connection to Aviasales autocomplete service"`
	}{}
	if _, err := flags.Parse(&opts); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("flags parsing error: %v\n", err)
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}
	if opts.Verbose {
		filter.SetMinLevel(logutils.LogLevel("DEBUG"))
	}
	logger := log.New(filter, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	logger.Printf("[INFO] Launching Application with: %+v", opts)

	boltDb, err := bolt.New(opts.DBPath)
	if err != nil {
		logger.Fatalf("[FATAL] Couldn't open bolt DB path: %v", err)
	}
	defer func() {
		if err = boltDb.Close(); err != nil {
			logger.Printf("[ERROR] Failed to close bolt DB: %v", err)
		}
	}()
	if err = boltDb.SetLogger(logger); err != nil {
		logger.Printf("[ERROR] Failed to set logger: %v", err)
	}

	avsProxy := proxy.New(opts.AVSAutocompleteURL, to.Duration(opts.AVSAutocompleteTimeout))
	if err = avsProxy.SetLogger(logger); err != nil {
		logger.Printf("[ERROR] Failed to set logger: %v", err)
	}

	httpServer := server.New(avsProxy, boltDb)
	if err = httpServer.SetLogger(logger); err != nil {
		logger.Printf("[ERROR] Failed to set logger: %v", err)
	}

	gr, appctx := errgroup.WithContext(context.Background())

	ErrCanceled := errors.New("Canceled")

	gr.Go(func() error {
		return httpServer.RunServer(appctx, opts.Listen, opts.Verbose)
	})

	gr.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		cusr := make(chan os.Signal, 1)
		signal.Notify(cusr, syscall.SIGUSR1)
		for {
			select {
			case <-appctx.Done():
				return nil
			case <-sigs:
				logger.Printf("[INFO] Caught stop signal. Exiting ...")
				return ErrCanceled
			case <-cusr:
				if filter.MinLevel == "DEBUG" {
					filter.SetMinLevel("INFO")
					logger.Printf("[INFO] Caught SIGUSR1 signal. Log level changed to INFO")
					continue
				}
				logger.Printf("[INFO] Caught SIGUSR1 signal. Log level changed to DEBUG")
				filter.SetMinLevel("DEBUG")
			}
		}
	})

	err = gr.Wait()
	if err != nil && err != ErrCanceled {
		logger.Fatal(err)
	}
}
