package margonem

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"scraper/database"
)

const baseUrl string = "https://public-api.margonem.pl"

type characterInfo struct {
	ID         string `json:"a"`
	Bucket     string `json:"c"`
	Nick       string `json:"n"`
	Profession string `json:"p"`
	Level      string `json:"l"`
	Dunno      string `json:"r"`
}

func CountMultipleWorldsInfoOnline(worlds []string) ([]database.CounterInsert, error) {
	var result []database.CounterInsert
	var wg sync.WaitGroup
	errorChan := make(chan error, len(worlds))

	for _, world := range worlds {
		wg.Add(1)
		go func(world string) {
			defer wg.Done()
			count, err := countWorldInfoOnline(world)
			if err != nil {
				errorChan <- fmt.Errorf("error for world %s: %v", world, err)
				return
			}
			counterInsert := database.CounterInsert{World: world, Count: count}
			result = append(result, counterInsert)
		}(world)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func countWorldInfoOnline(world string) (int, error) {
	apiUrl := composeUrlToWorldInfoOnline(world)
	res, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var data []characterInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func composeUrlToWorldInfoOnline(world string) string {
	return fmt.Sprintf("%s/info/online/%s.json", baseUrl, world)
}
