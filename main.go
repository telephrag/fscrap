package main

import (
	"fmt"
	"log"
	"net/http"
	"softbuilding/config"
	"softbuilding/fabrikant"

	"github.com/PuerkitoBio/goquery"
)

var query = "Ремонт"

func init() {} // parse flags for query

func main() {

	// vary values for page, on_page, query etc.
	url := fmt.Sprintf(
		"https://www.fabrikant.ru/trades/procedure/search/?type=0&procedure_stage=0&currency=0&date_type=date_publication&ensure=all&count_on_page=40&order_direction=1&query=%s&page=1",
		query,
	)
	req, _ := http.NewRequest("GET", url, nil)
	// req.AddCookie(resp.Cookies()[0])
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	orders := make([]*fabrikant.Order, config.ON_PAGE)[:0]
	// fmt.Println(doc.Find(".marketplace-list").Find(".innerGrid").Text())
	doc.Find(".marketplace-list").
		Find(".innerGrid").
		Each(func(i int, s *goquery.Selection) {
			orders = append(orders, fabrikant.NewOrderFromInnerGridSelection(s))
		})

	for _, o := range orders {
		fmt.Println(o.UID, o.PublicationTimestamp)
	}
}
