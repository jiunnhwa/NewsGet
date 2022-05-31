package model

import (
	"time"
)

//Feed is data model of NewsItem
type Feed struct {
	Author        string
	Title         string
	Description   string
	SourceID      string
	SourceName    string
	URL           string
	PublishedTime time.Time //mysql: datetime
}

//NewFeedItem with default values, else NULL issues
func NewFeedItem() *Feed {
	return &Feed{Title: "", Author: "", Description: "", URL: "", PublishedTime: time.Time{}}
}
