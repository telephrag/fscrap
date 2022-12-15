package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func main() {
	query := "Ремонт"

	// vary values for page, on_page, query etc.
	url := fmt.Sprintf(
		"https://www.fabrikant.ru/trades/procedure/search/?query=%s&type=0&org_type=org&currency=0&date_type=date_publication&ensure=all&okpd2_embedded=1&okdp_embedded=1&type_hash=1671021652&count_on_page=40&page=1",
		query,
	)
	req, _ := http.NewRequest("GET", url, nil)
	// req.AddCookie(resp.Cookies()[0])
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// body := make([]byte, resp.ContentLength)
	// _, err = resp.Body.Read(body)
	// if err != nil && err != io.EOF {
	// 	log.Fatal(err)
	// }

	// var content any
	// err = json.Unmarshal(body, &content)
	// if err != nil {
	// 	// log.Fatal(err)
	// }

	// fmt.Println(string(resp.Body.Read(body)))

	buff := bufio.NewScanner(resp.Body)
	buff.Split(bufio.ScanLines)
	for buff.Scan() {
		fmt.Println(buff.Text())
	}
	fmt.Println(buff.Text())
}
