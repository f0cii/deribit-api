package models

type Announcement struct {
	Title           string `json:"title"`
	PublicationTime int64  `json:"publication_time"`
	Important       bool   `json:"important"`
	ID              int64  `json:"id"`
	Body            string `json:"body"`
}
