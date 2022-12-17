package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"softbuilding/config"
	"softbuilding/fabrikant"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ltime | log.Ldate))
}

func main() {
	// doesn' work properly if placed inside init() for some reason

	query := flag.String("q", "", "Query to search orders with. Use 'your query here' for multiword query.")
	withUID := flag.Bool("uid", false, "Set to print uids of orders.")
	withType := flag.Bool("type", false, "Set to print types of orders.")
	withTitle := flag.Bool("title", false, "Set to print titles of orders.")
	withTimestamp := flag.Bool("timestamp", false, "Set to print publication timestamps of orders.")
	flag.Parse()

	// fmt.Println(*query)
	// fmt.Println(*withUID)
	// fmt.Println(*withType)
	// fmt.Println(*withTitle)
	// fmt.Println(*withTimestamp)

	queryEscaped := url.QueryEscape(*query)

	url := fmt.Sprintf(
		"https://www.fabrikant.ru/trades/procedure/search/?type=0&procedure_stage=0&price_from=&price_to=&currency=0&date_type=date_publication&date_from=&date_to=&ensure=all&count_on_page=40&order_direction=1&type_hash=1561441166&query=%s",
		queryEscaped,
	)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		// s := bufio.NewScanner(resp.Body)
		// for s.Scan() {
		// 	fmt.Println(s.Text())
		// }
		log.Fatal("Error making http request: ", resp.StatusCode)
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

	printUid := func(o *fabrikant.Order) {
		fmt.Print(o.UID, " | ")
	}
	printType := func(o *fabrikant.Order) {
		fmt.Print(o.Type, " | ")
	}
	printTitle := func(o *fabrikant.Order) {
		fmt.Print(o.Title, " | ")
	}
	printTimestamp := func(o *fabrikant.Order) {
		fmt.Print(o.PublicationTimestamp, " | ")
	}

	printers := []func(o *fabrikant.Order){}
	if *withUID {
		printers = append(printers, printUid)
	}
	if *withType {
		printers = append(printers, printType)
	}
	if *withTitle {
		printers = append(printers, printTitle)
	}
	if *withTimestamp {
		printers = append(printers, printTimestamp)
	}

	if len(printers) == 0 {
		return
	}

	for _, o := range orders {
		for _, p := range printers {
			p(o)
		}
		fmt.Println()
	}
}
