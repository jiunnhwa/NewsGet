package newsapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/endeveit/guesslanguage"

	"newsget/manager/job"
	"newsget/model"
	webclient "newsget/service/http/client"
	myurl "newsget/service/http/url"

	"database/sql"
	"newsget/model/sql/mysql"

	mylogger "newsget/service/logger"
)

func NewJobAPI() {
	job := job.NewJob("123")
	defer func() {
		job.End()
		job = nil
	}()

	job.NewsData.AsBytes = (webclient.Fetch("GET", "https://mocki.io/v1/e750d778-4861-498e-b00e-213314f799d6", ""))
	job.NewsData.AsString = string(job.NewsData.AsBytes)
	job.Feeds = ProcessData(job.NewsData.AsBytes)
	fmt.Println(job.Feeds)
}

//Run the gets, concat the results
func Run(db *sql.DB, urls []string) (feeds []model.Feed) {
	//Download and Process
	for _, url := range urls {
		mylogger.Trace.Println("Run:", url)
		f := ProcessData(Download(url))

		feeds = append(feeds, f...)
	}
	//Insert to DB
	Upload(db, feeds)
	//Return
	return feeds
}

//Download from url as bytes
func Download(url string) []byte {
	return webclient.Fetch("GET", url, "")
}

//ProcessData transforms bytes to model
func ProcessData(bytes []byte) []model.Feed {
	Feeds := []model.Feed{}
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		mylogger.Trace.Println("ProcessData:\t", err.Error())
		return Feeds
	}
	body := string(bytes)
	if InStr(body, []string{"error", "rateLimited"}) {
		//{"status":"error","code":"rateLimited","message":"You have made too many requests recently. Developer accounts are limited to 100 requests over a 24 hour period (50 requests available every 12 hours). Please upgrade to a paid plan if you need more requests."}
		mylogger.Trace.Println("ProcessData:\t", body)
		return Feeds
	}
	if result == nil {
		mylogger.Trace.Println("ProcessData:\t", "result is nil")
		return Feeds
	}
	articles := result["articles"]
	if articles == nil {
		mylogger.Trace.Println("ProcessData:\t", "articles is nil")
		return Feeds
	}
	obj := articles.([]interface{})
	mylogger.Trace.Println("\tArticles#:", len(obj))
	for _, article := range obj {

		// decode
		source := article.(map[string]interface{})["source"]
		sourceid := source.(map[string]interface{})["id"]
		sourcename := source.(map[string]interface{})["name"]
		author := article.(map[string]interface{})["author"]           //can be null
		title := article.(map[string]interface{})["title"].(string)    //assert
		description := article.(map[string]interface{})["description"] //can be null
		url := article.(map[string]interface{})["url"]
		publishedAt := article.(map[string]interface{})["publishedAt"]

		//Data cleanup
		u := ToString(url)
		sn := ToString(sourcename)
		if !IsGuessLangEng(title) || IsNonEnglish(u) || IsPaywall(u) || IsBadSourceName(sn, u) {
			//Do nothing
		} else {
			Feeds = append(Feeds, model.Feed{Author: ToString(author), Title: (title), Description: ToString(description), SourceID: ToString(sourceid), SourceName: ToString(sourcename), URL: ToString(url), PublishedTime: ToTime(time.RFC3339, ToString(publishedAt))}) //of each article
		}

	}
	mylogger.Trace.Println("\tFeeds#:", len(Feeds))
	return Feeds
}

//Upload rows
func Upload(db *sql.DB, feeds []model.Feed) (Affected int64) {
	for _, v := range feeds {
		Affected, _ = mysql.AddNewsItem(db, v)
		//Error 1062: Duplicate entry '...' for key 'PRIMARY' recordsAffected -1
	}
	return Affected
}

//IsGuessLangEng checks whether it is english
func IsGuessLangEng(text string) bool {
	lang, err := guesslanguage.Guess(text)
	if err != nil {
		mylogger.Trace.Println("IsGuessLangEng:\t", err, lang, text)
	}
	mylogger.Trace.Println("IsGuessLangEng:\t", lang == "en", lang, text)
	return (lang == "en")
}

//IsNonEnglish denys non-english news site
func IsNonEnglish(url string) bool {
	//use deny all others rule
	if InStr(myurl.GetHostSuffix(url), []string{"com", "co", "ca" /*cbc.ca,canada*/}) {
		//non-english .com
		if InStr(url, []string{"russian.rt.com", "handelsblatt.com", "globo.com", "infobae.com",
			"cnbeta.com", "wallstreetcn.com", "cnwest.com", "hk01.com", "xataka.com",
			"journalducoin.com", "lavanguardia.com", "elperiodico.com", "hobbyconsolas.com", "sa.investing.com",
			"cnnespanol.cnn.com", "eluniverso.com", "culturaocio.com", "efe.com", "br.investing.com", "marca.com", "almasryalyoum.com",
			"diepresse.com", "elespanol.com", "actualidad.rt.com", "mdzol.com", "ambito.com", "lesnumeriques.com",
			"antena3.com",
		}) {
			mylogger.Trace.Println("IsNonEnglish:\t", url)
			return true
		}
		//valid .com
		return false
	}
	mylogger.Trace.Println("IsNonEnglish:\t", url)
	return true
}

//IsBadSourceName - bad listing
func IsBadSourceName(sourcename, url string) bool {
	if InStr(sourcename, []string{"Lenta", "RBC"}) {
		mylogger.Trace.Println("IsBadSourceName:\t", url)
		return true
	}
	return false
}

//IsPaywall - paid sites
func IsPaywall(url string) bool {
	if InStr(url, []string{"ft.com", "express.co.uk", "fortune.com"}) {
		mylogger.Trace.Println("IsPaywall:\t", url)
		return true
	}
	return false
}

//converts interface to string value
func ToString(x interface{}) string {
	//panic: interface conversion: interface {} is nil, not string
	if _, ok := x.(string); !ok {
		return ""
	}
	return fmt.Sprintf("%v", x) //nil interface returns "\u003cnil\u003e"
}

//Returns time value, or timezero on err
func ToTime(layout, timestr string) time.Time {
	//const layout = time.RFC3339 //"2006-01-02T15:04:05Z"  newsapi layout
	// Calling Parse() method with its parameters
	tm, err := time.Parse(layout, timestr)
	if err != nil {
		// time.Time{} returns January 1, year 1, 00:00:00.000000000 UTC
		// which according to the source code is the zero value for time.Time
		// https://golang.org/src/time/time.go#L23
		return time.Time{}

	}
	return tm
}

//InStr check if list contains item
func InStr(str string, list []string) bool {
	ustr := strings.ToUpper(strings.TrimSpace(str))
	for _, v := range list {
		if strings.Contains(ustr, strings.ToUpper(strings.TrimSpace(v))) {
			return true
		}
	}
	return false
}
