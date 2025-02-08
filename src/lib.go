package src

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func Main() {
	flags := parseFlags()

	if flags.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// load or create checkpoint
	log.Info("Loading or creating checkpoint file")
	checkpoint, err := loadCheckpoint(flags)
	if err != nil {
		log.Error("Failed to load checkpoint")
		log.Fatal(err)
	}

	log.Debug("Initializing Rod")
	u := launcher.New().Headless(true).Set("--disable-gpu").Set("--no-sandbox").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// create page
	page := browser.MustPage("")
	defer page.MustClose()

	// setup error recovery
	page.MustHandleDialog()

	// handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// create error channel
	errChan := make(chan error, 1)

	// create context with cancelation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run scraper in goroutine
	go func() {
		switch flags.Mode {
		case MODE_ATTENTE:
			log.Info("Scraping liste d'attente")
			errChan <- scrapeWaitingList(ctx, page, checkpoint, flags)
		case MODE_OFFICINE:
			log.Info("Scraping officines")
			errChan <- scrapeOfficines(ctx, page, checkpoint, flags)
		}
	}()

	// wait  for either completion or interruption
	select {
	case err := <-errChan:
		if err != nil {
			log.Error("Error while scraping waiting list")
			log.Error("You can resume the scraping by running the program again")
			log.Fatal(err)
		}

		log.Info("Writing results")
		if err := writeResults(checkpoint, flags); err != nil {
			log.Error("Failed to write results")
			log.Fatal(err)
		}

		// cleanup checkpoint file on success
		if err := os.Remove(checkpointFileName(flags)); err != nil {
			log.Error("Failed to remove checkpoint file")
			log.Error(err)
		}
	case sig := <-sigChan:
		log.Infof("Received signal %s, saving checkpoint and exiting", sig)
		if err := saveCheckpoint(checkpoint, flags); err != nil {
			log.Error("Failed to save checkpoint")
		}
		os.Exit(1)
	}

}
