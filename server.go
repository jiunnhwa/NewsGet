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
	//viewData := &ViewData{PageTitle: "Home", Msg: "Welcome to NewsGetter."}
	//ViewHome.SetViewData(viewData).ServeTemplate(w, r)
}

//apijobs returns a list of News
//Order by RID Desc(lastest on top)
func apiNews(w http.ResponseWriter, r *http.Request) {
	// if only one expected
	//param1 := r.URL.Query().Get("param1")
	//if param1 != "" {
	// ... process it, will be the first (only) if multiple were given
	// note: if they pass in like ?param1=&param2= param1 will also be "" :|
	//}
	// if multiples possible, or to process empty values like param1 in
	// ?param1=&param2=something
	//param1s := r.URL.Query()["param1"]
	//if len(param1s) > 0 {}
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
