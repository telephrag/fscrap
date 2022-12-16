package fabrikant

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Order struct {
	UID                  string
	Type                 string
	Title                string
	URL                  string
	PublicationTimestamp time.Time
}

func NewOrderFromInnerGridSelection(s *goquery.Selection) *Order {
	o := &Order{}

	inner := s.Find(".marketplace-unit")

	// uid and type
	uidFull := inner.Find(".marketplace-unit__info__name")
	uidData := uidFull.Nodes[0].FirstChild.NextSibling.FirstChild.Data
	uidSplit := strings.Split(uidData, "№")
	if len(uidSplit) < 2 || len(uidSplit) > 2 {
		log.Fatalf("Failed to parse number and type of order: \"%s\"\n", uidFull)
	}
	o.Type, o.UID = strings.TrimSpace(uidSplit[0]), strings.TrimSpace(uidSplit[1])

	// title and url
	title := inner.Find(".marketplace-unit__cut-wrap").
		Find(".marketplace-unit__title").
		Find(".text")
	o.Title = strings.TrimSpace(title.Text())
	var ok bool
	if o.URL, ok = title.Attr("href"); !ok {
		log.Fatalf("Failed to parse URL of order: \"%s\"\n", title.Text())
	}

	// publication time
	pubDate := inner.Find(".marketplace-unit__state__wrap").
		Find(".marketplace-unit__state").First().Children()

	if len(pubDate.Nodes) < 3 {
		log.Fatalf("Failed to parse date and time, not enough nodes: \"%s\"\n", pubDate.Text())
	}
	dt := pubDate.Nodes[1].FirstChild.Data + " " + pubDate.Nodes[2].FirstChild.Data
	var err error
	if o.PublicationTimestamp, err = parseTimestamp(dt); err != nil {
		log.Fatal(err)
	}

	return o
}

func parseTimestamp(ts string) (time.Time, error) {
	if len(ts) < 10 { // 10 is minimum length of timestamp on fabrikant.ru
		return time.Time{}, fmt.Errorf("too short timestamp: \"%s\"", ts)
	}

	// parsing month token
	sp := strings.Split(ts, " ")
	if len(sp) != 4 {
		return time.Time{}, fmt.Errorf("not enough tokens in timestamp string's split: \"%s\"", ts)
	}
	m := sp[1]

	switch {
	case m == "янв":
		m = "01"
	case m == "фев":
		m = "02"
	case m == "мар":
		m = "03"
	case m == "апр":
		m = "04"
	case m == "май":
		m = "05"
	case m == "июн":
		m = "06"
	case m == "июл":
		m = "07"
	case m == "авг":
		m = "08"
	case m == "сен":
		m = "09"
	case m == "окт":
		m = "10"
	case m == "ноя":
		m = "11"
	case m == "дек":
		m = "12"
	default:
		return time.Time{}, fmt.Errorf("no month token was found in timestamp string's split: \"%s\"", ts)
	}

	return time.ParseInLocation(
		"2/01/2006 15:04", sp[0]+"/"+m+"/"+sp[2]+" "+sp[3],
		time.FixedZone("UTC", 10800),
	)
}
