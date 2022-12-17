package fabrikant

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ( // errors
	ErrIDAndTypeParse = errors.New("failed to parse id and type")
	ErrURLParse       = errors.New("failed to parse url")
	ErrDateTime       = errors.New("failed to parse date and time, not enough nodes")

	ErrTimestampNotEnoughTokens = errors.New("not enough tokens in timestamp string's split")
	ErrTimestampNoMonth         = errors.New("no month token")
)

type Order struct {
	UID                  string
	Type                 string
	Title                string
	URL                  string
	PublicationTimestamp time.Time
}

func NewOrderFromInnerGridSelection(s *goquery.Selection) (*Order, error) {
	o := &Order{}

	inner := s.Find(".marketplace-unit")

	// uid and type
	uidFull := inner.Find(".marketplace-unit__info__name")
	uidData := uidFull.Nodes[0].FirstChild.NextSibling.FirstChild.Data
	uidSplit := strings.Split(uidData, "№")
	if len(uidSplit) < 2 || len(uidSplit) > 2 {
		return nil, fmt.Errorf("%w: \"%s\"", ErrIDAndTypeParse, uidData)
	}
	o.Type, o.UID = strings.TrimSpace(uidSplit[0]), strings.TrimSpace(uidSplit[1])

	// title and url
	title := inner.Find(".marketplace-unit__cut-wrap").
		Find(".marketplace-unit__title").
		Find(".text")
	o.Title = strings.TrimSpace(title.Text())
	var ok bool
	if o.URL, ok = title.Attr("href"); !ok {
		return nil, fmt.Errorf("%w: \"%s\"", ErrURLParse, title.Text())
	}

	// publication time
	pubDate := inner.Find(".marketplace-unit__state__wrap").
		Find(".marketplace-unit__state").First().Children()

	if len(pubDate.Nodes) < 3 {
		return nil, fmt.Errorf("%w: \"%s\"", ErrDateTime, pubDate.Text())
	}
	dt := pubDate.Nodes[1].FirstChild.Data + " " + pubDate.Nodes[2].FirstChild.Data
	var err error
	if o.PublicationTimestamp, err = parseTimestamp(dt); err != nil {
		return nil, fmt.Errorf("%w: \"%s\"", err, dt)
	}

	return o, nil
}

func parseTimestamp(ts string) (time.Time, error) {

	// parsing month token
	sp := strings.Split(ts, " ")
	if len(sp) != 4 {
		return time.Time{}, ErrTimestampNotEnoughTokens
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
		return time.Time{}, ErrTimestampNoMonth
	}

	return time.ParseInLocation(
		"2/01/2006 15:04", sp[0]+"/"+m+"/"+sp[2]+" "+sp[3],
		time.FixedZone("UTC", 10800),
	)
}
