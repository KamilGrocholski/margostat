package margonem

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"scraper/database"
)

func CountMultipleWorldsInfoOnline() ([]database.CounterInsert, error) {
	url := "https://www.margonem.pl/stats"
	userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	popups := doc.Find("div.light-brown-box.news-container.no-footer[class$=\"-popup\"]")

	servers := make([]database.CounterInsert, popups.Length())

	popups.Each(func(popupIndex int, popup *goquery.Selection) {
		class, _ := popup.Attr("class")
		classList := strings.Split(class, " ")
		var serverClass string
		for _, class := range classList {
			serverClass = class
		}
		serverName := strings.ReplaceAll(serverClass, "-popup", "")

		namesList := doc.Find(fmt.Sprintf("div.%s-popup .statistics-rank", serverName))

		count := namesList.Length()

		servers[popupIndex] = database.CounterInsert{
			World: serverName,
			Count: count,
		}
	})

	return servers, nil
}
