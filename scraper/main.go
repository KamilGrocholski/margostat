package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"scraper/config"
	"scraper/database"
	"scraper/logger"
	"scraper/margonem"
	"scraper/ssg"
)

func main() {
	logger.Init()
	config.Init()

	c := cron.New()

	cronSpec := fmt.Sprintf("*/%d * * * *", config.SCRAPER_INTERVAL_IN_MINUTES)

	c.AddFunc(cronSpec, func() {
		db, err := database.Connect()
		if err != nil {
			logger.Logger.Fatalf("Failed connecting to db, %v", err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				logger.Logger.Printf("Error while closing db connection, %v", err)
			} else {
				logger.Logger.Print("Closed db connection")
			}
		}()
		if err := database.Init(db); err != nil {
			logger.Logger.Fatalf("Failed init on db, %v", err)
		}

		logger.Logger.Printf("Start scraping")

		for attempt := 1; attempt <= int(config.SCRAPER_MAX_ATTEMPTS); attempt++ {
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
		GeneratedAt:   time.Now(),
	}
	allGenErr := ssg.GenerateAndWriteHtmlPageFileToStatic("index", &allWorldsPageData)
	if allGenErr != nil {
		logger.Logger.Printf("Failed index page generation %v", allGenErr)
	}

	worldsGeneratedCount := 0

	for _, w := range worlds {
		timeline, err := database.GetWorldTimeline(db, w)
		if err != nil {
			logger.Logger.Fatalf("error getting timeline for world %s, %v", w, err)
		}
		data := ssg.WorldPageData{
			SelectedWorld: w,
			CountResults:  timeline,
			Worlds:        worlds,
			GeneratedAt:   time.Now(),
		}
		genErr := ssg.GenerateAndWriteHtmlPageFileToStatic(w, &data)
		if genErr != nil {
			logger.Logger.Printf("Failed world %s page generation %v", w, genErr)
		} else {
			worldsGeneratedCount++
		}
	}

	logger.Logger.Printf("Generated %d/%d worlds", worldsGeneratedCount, len(worlds))

	return nil
}
