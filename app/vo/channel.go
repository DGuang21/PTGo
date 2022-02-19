package vo

import "encoding/xml"

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string    `xml:",chardata"`
		Title string    `xml:"title"`
		Link  string    `xml:"link"`
		Item  []RssItem `xml:"item"`
	} `xml:"channel"`
}

type RssItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Category    struct {
		Text   string `xml:",chardata"`
		Domain string `xml:"domain,attr"`
	} `xml:"category"`
	Comments  string `xml:"comments"`
	Enclosure struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	Guid struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	PubDate string `xml:"pubDate"`
}
