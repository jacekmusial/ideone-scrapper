package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("ideone.com"))

	c.OnHTML("strong", func(htmlElement *colly.HTMLElement) {
		link := htmlElement.Attr("href")
		fmt.Printf("Found link: %q -> %s\n", htmlElement.Text, link)
	})

	c.Visit("https://ideone.com/recent")

}
