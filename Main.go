package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)


const (
	DB_HOST = "tcp(192.168.1.65:3306)"
	DB_NAME = "ideone"
	DB_USER = /*"root"*/ "root"
	DB_PASS = /*""*/ "nanomader#!$!%("
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
	//jdbc:mariadb://192.168.1.65:3306/ideone
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	checkErr(err)

	split := strings.Split(links, ";")
	split = split[1:50]
	for i, k := range split {
		fmt.Println("----------")
		var url = "https://ideone.com/plain/" + k[1:]

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			html, err := ioutil.ReadAll(response.Body)
			var txt string = string(html)
			fmt.Println(len(html))
			//	fullurl VARCHAR(60),
			//	codedate VARCHAR(60),
			//	codekey VARCHAR(30),
			//	size BIGINT,
			//	codelines BIGINT,
			//	language VARCHAR(30),
			//	status VARCHAR(30),
			//  txt LONGTEXT,
			stmt, err := db.Prepare("INSERT INTO IE (fullurl, codedate, codekey, size, codelines, language, " +
				"status, txt) VALUES (?,?,?,?,?,?,?,?)")
			checkErr(err)

			currentTime := time.Now()
			res, err := stmt.Exec(url, currentTime.Format("2006-01-02 15:04:05"), k[1:], len(html),
				strings.Count(txt, "\n"), "language", result[i], txt)
			if res == nil {

			}
			checkErr(err)
		}

		fmt.Println(result[i], ", ", url)
		fmt.Println("____")
	}

}

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR!")
		log.Fatal(err)
		panic(err)
	}
}