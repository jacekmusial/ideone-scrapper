package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)


const (
	DB_HOST = "tcp(localhost:3306)"
	DB_NAME = "ideone"
	DB_USER = /*"root"*/ "root"
	DB_PASS = /*""*/ "password"
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
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	var db *sql.DB
	db, err := sql.Open("mysql", dsn)
	checkErr(err)
	memStorage := storage.NewMemoryStorage()
	s := scheduler.New(memStorage)
	fmt.Println("Im good at school three huna worldwide")
	if _, err := s.RunEvery(40 * time.Second, scrapIdeone, db); err != nil {
		log.Fatal(err)
	}
	s.Start()
	s.Wait()
}

func scrapIdeone(db *sql.DB) {
	links, result := getRecentLinks()
	//jdbc:mariadb://192.168.1.65:3306/ideone

	split := strings.Split(links, ";")
	split = split[1:50]

	results := strings.Split(result, ";")
	results = results[1:50]

	for i, k := range split {
		fmt.Println("----------")
		var url = "https://ideone.com/plain/" + k[1:]

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			html, err := ioutil.ReadAll(response.Body)
			print(response.Header.Get("Content-Disposition"))
			get := response.Header.Get("Content-Disposition")
			fmt.Println("get: ", get)
			start := strings.Index(get, "=\"")
			fmt.Println(start)
			filename := get[(start + 1):]
			language := filename[strings.Index(filename, ".")+1 : len(filename)-1]

			var txt = string(html)
			fmt.Println(len(html))
			stmt, err := db.Prepare("INSERT INTO IE (fullurl, codedate, codekey, size, codelines, language, " +
				"status, txt) VALUES (?,?,?,?,?,?,?,?)")
			checkErr(err)

			currentTime := time.Now()
			_, err = stmt.Exec(url, currentTime.Format("2006-01-02 15:04:05"), k[1:], len(html),
				strings.Count(txt, "\n"), language, results[i], txt)
			if err != nil {
				me, _ := err.(*mysql.MySQLError)
				if me.Number == 1062 {
					fmt.Println("It already exists in a database. ", err)
				} else {
					checkErr(err)
				}
			}
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