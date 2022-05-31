package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"newsget/jutil/v1/http/response"
	"newsget/service/text"
	"strconv"
	"time"
)

type Feed struct {
	Title string
	URL   string

	Description  string
	Note         string
	CategoryPath string //Golang/CheatSheet/...
	Tags         []string

	RID         int //RecordID:LastID
	CreatedDate time.Time
	CreatedBy   string //owner of the collection
}

//ServeRoutes defines the routing
func ServeRoutes() error {

	//VIEWS
	http.HandleFunc("/", home)

	//API
	http.HandleFunc("/api/v1/newsget/news", apiNews)
	fmt.Println("Newsget listening...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}

//Views are globals, loaded once, and reused.
//var ViewHome = NewView([]string{filepath.Clean(filepath.Join(tplDir, "home.gohtml")), filepath.Join(tplDir, "base.gohtml")})

//handles home page, updates the view data and serve
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to NewsGetter.")
}

//apiNews handles the call
func apiNews(w http.ResponseWriter, r *http.Request) {
	limit := -1
	if n, err := strconv.Atoi(text.ToUpper(r.URL.Query().Get("limit"))); err == nil {
		limit = n
	}

	if r.Method == http.MethodGet {
		feeds, _ := GetNewsLatest(DB, limit)
		fmt.Println(feeds)
		response.AsJSON(w, 200, feeds)
		return
	}
	response.AsJSON(w, http.StatusBadRequest, "Invalid Action.")
}

//GetNewsLatest returns a result set
//order by PublishedTime DESC with LIMIT
func GetNewsLatest(db *sql.DB, limit int) (*[]Feed, error) {
	sql := ""
	if limit > 0 {
		sql = fmt.Sprintf("SELECT  Title, URL FROM goschool.news ORDER BY PublishedTime DESC LIMIT %d", limit)
	} else {
		sql = "SELECT  Title, URL FROM goschool.news ORDER BY PublishedTime DESC "
	}
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Feed
	for rows.Next() {
		item := Feed{}
		if err := rows.Scan(&item.Title, &item.URL); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return &result, nil
}
