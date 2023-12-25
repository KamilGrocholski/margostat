package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"

	"scraper/database"
	"scraper/logger"
	"scraper/margonem"
	"scraper/ssg"
)

func main() {
	scraperIntervalInMinutes, err := strconv.ParseInt(os.Getenv("SCRAPER_INTERVAL_IN_MINUTES"), 10, 32)
	if err != nil {
		logger.Logger.Fatalf("parsing interval env, %v", err)
	}

	scraperMaxAttempts, err := strconv.ParseInt(os.Getenv("SCRAPER_MAX_ATTEMPTS"), 10, 32)
	if err != nil {
		logger.Logger.Fatalf("parsing max attempts env, %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		logger.Logger.Printf("failed connecting to db, %v", err)
	}
	defer db.Close()
	if err := database.Init(db); err != nil {
		logger.Logger.Printf("failed init on db, %v", err)
	}

	c := cron.New()

	cronSpec := fmt.Sprintf("*/%d * * * *", scraperIntervalInMinutes)
	c.AddFunc(cronSpec, func() {
		logger.Init()
		logger.Logger.Printf("Start scraping")

		for attempt := 1; attempt <= int(scraperMaxAttempts); attempt++ {
			err := scrap(db)
			if err == nil {
				logger.Logger.Printf("Success scraping (attempt %d): %v\n", attempt, err)
				break
			}

			logger.Logger.Printf("Error scraping (attempt %d): %v\n", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		revalidate(db)
	})

	c.Start()

	select {}
}

func scrap(db *sql.DB) error {
	counters, counterErr := margonem.CountMultipleWorldsInfoOnline(margonem.MARGONEM_WORLD_NAMES)
	if counterErr != nil {
		return counterErr
	}

	insertErr := database.InsertMultipleCounters(db, counters)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func revalidate(db *sql.DB) error {
	worlds, err := database.GetAllWorldNames(db)
	if err != nil {
		return err
	}

	allWorldsTimeline, err := database.GetGlobalTimeline(db)
	if err != nil {
		return err
	}
	allWorldsPageData := ssg.WorldPageData{
		SelectedWorld: margonem.MARGONEM_GLOBAL_NAME,
		Worlds:        worlds,
		CountResults:  allWorldsTimeline,
	}
	allGenErr := ssg.GenerateAndWriteHtmlPageFileToStatic("index", &allWorldsPageData)
	if allGenErr != nil {
		logger.Logger.Printf("index not generated %v", allGenErr)
	}

	for _, w := range worlds {
		timeline, err := database.GetWorldTimeline(db, w)
		if err != nil {
			logger.Logger.Fatalf("world %s not timeline, %v", w, err)
		}
		data := ssg.WorldPageData{
			SelectedWorld: w,
			CountResults:  timeline,
			Worlds:        worlds,
		}
		genErr := ssg.GenerateAndWriteHtmlPageFileToStatic(w, &data)
		if genErr != nil {
			logger.Logger.Printf("world %s not generated %v", w, genErr)
		}
	}

	return nil
}
