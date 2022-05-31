package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	"newsget/manager/job"
	webclient "newsget/service/http/client"
	"newsget/service/http/url"
	mylogger "newsget/service/logger"
	"path"
	"time"

	"newsget/provider/newsapi"

	_ "github.com/go-sql-driver/mysql"
)

const APIKEY = "YOUR_API_KEY"

var URLs []string

//Results represents the decoded json fields
type Results struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

var DB *sql.DB

func main() {
	mylogger.LogDir = "logs"
	mylogger.InitLogger()
	mylogger.GoStartLogRotator()

	DB = OpenDB()
	defer CloseDB()

	urlAsia := LoadURLs("urlAsia.txt")
	urlUS := LoadURLs("urlUS.txt")
	urlTwice := LoadURLs("urlTwice.txt")

	//Time based download based on the Region Active Hours
	go func() {
		ticker := time.NewTicker(time.Minute)
		for ; true; <-ticker.C {
			localHH, localNN := time.Now().Local().Hour(), time.Now().Local().Minute()
			if localNN == 15 {
				//Asia
				if localHH == 7 || localHH == 9 || localHH == 10 || localHH == 13 || localHH == 17 || localHH == 21 {
					DoNewsAPI(urlAsia)
				}

				//US
				if localHH == 20 || localHH == 21 || localHH == 22 || localHH == 0 || localHH == 2 || localHH == 4 || localHH == 6 {
					DoNewsAPI(urlUS)
				}

				//Twice a day
				if localHH == 11 || localHH == 18 {
					DoNewsAPI(urlTwice)
				}
			}
		}
	}()

	ServeRoutes()

	fmt.Println("Press Enter To Exit")
	fmt.Scanln()
}

func DoNewsAPI(urls []string) {
	job := job.NewJob("newsapi")
	mylogger.Trace.Println("START:", job.ID)
	defer func() {
		mylogger.Trace.Println("END:", job.ID)
		job.End()
		job = nil
	}()
	newsapi.Run(DB, urls)
}

func GetNews(urls []string, interval time.Duration) {
	NextURL := func(list []string) func() string {
		len := len(list)
		currNum := -1
		return func() string {
			currNum += 1
			if currNum >= len {
				currNum = 0
			}
			return list[currNum]
		}
	}(urls)

	ticker := time.NewTicker(interval)
	for ; true; <-ticker.C {
		u := NextURL()
		fname := path.Join("inbox", url.GetHostName(u)+"-"+"data.txt")
		fmt.Println(fname, u)
		WriteData(fname, (webclient.Fetch("GET", u, "")))
	}
}

//LoadURLs from store
func LoadURLs(ffname string) []string {
	links := []string{}
	if data, err := ioutil.ReadFile(ffname); err == nil {
		apis := strings.Split(string(data), "\n")
		for _, v := range apis {
			if len(v) > 0 {
				links = append(links, v) //url
			}
		}
		fmt.Println("LoadURLs", links)
		return links

	}

	return links
}
