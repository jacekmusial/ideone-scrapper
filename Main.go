package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getRecentLinks() (string, string) {
	c := colly.NewCollector(colly.AllowedDomains("ideone.com"))

	var ret string
	var result string

	c.OnHTML("strong", func(htmlElement *colly.HTMLElement) {
		ret += ";" + htmlElement.Text
	})

	c.OnHTML("span .info", func(htmlElement *colly.HTMLElement) {
		result += ";" + htmlElement.Text
	})

	c.Visit("https://ideone.com/recent")

	return ret, result
}

func main() {
	links, result := getRecentLinks()
	fmt.Print(result)

	split := strings.Split(links, ";")
	split = split[1:49]
	for _, k := range split {
		fmt.Printf("https://ideone.com/plain/%s \n", k[1:])
		var url = "https://ideone.com/plain/" + k[1:]

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			_, err := io.Copy(os.Stdout, response.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//jdbc:mariadb://192.168.1.65:3306/ideone
	//db, err := sql.Open("mysql", "jdbc:mariadb://192.168.1.65:3306/ideone")
	//checkErr(err)


}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}