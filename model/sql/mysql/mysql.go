package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"newsget/model"

	mylogger "newsget/service/logger"
)

//GetNews from table and return as Model
func GetNews(db *sql.DB) ([]model.Feed, error) {
	rows, err := db.Query("Select Title, Description ,URL, PublishedTime FROM goschool.news LIMIT 1000")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Feeds := []model.Feed{}
	for rows.Next() {
		var f model.Feed
		if err := rows.Scan(&f.Title, &f.Description, &f.URL, &f.PublishedTime); err != nil {
			return nil, err
		}
		Feeds = append(Feeds, model.Feed{Title: f.Title, Description: f.Description, URL: f.URL, PublishedTime: f.PublishedTime})
	}

	return Feeds, nil
}

//AddNewsItem to table
func AddNewsItem(db *sql.DB, f model.Feed) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO goschool.news (`Title`, `Author`, `Description`, `SourceName`, `URL`, `PublishedTime`) VALUES (?,?,?,?,?,?);",
		f.Title, f.Author, f.Description, f.SourceName, f.URL, f.PublishedTime)

	if err != nil {
		//Error 1062: Duplicate entry 'YODA123' for key 'PRIMARY'
		//Error 1406: Data too long for column 'Author' at row 1
		fmt.Println("AddNewsItem>", f.Title, "\n", err)
		mylogger.Error.Println("AddNewsItem>", f.Title, "\n", err)
		return -1, err
		//log.Fatal(err)
	}
	/*lastId*/
	_, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected, nil
}

//*************************************************************
// SQL Sanitizer
//*************************************************************
//StringEscape escapes the input string accordingly
func StringEscape(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	str = strings.Replace(str, ";", "\\;", -1)
	str = strings.Replace(str, "\"", "\\'", -1)
	return str
}

//SQLReject guards against SQL injection
func SQLReject(str string) error {
	fields := strings.Fields(str)
	if len(fields) > 1 {
		return errors.New("wrong number of fields in input")
	}
	return nil
}
