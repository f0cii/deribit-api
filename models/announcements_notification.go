package models

type AnnouncementsNotification struct {
	Action    string `json:"action"`
	Body      string `json:"body"`
	Date      int64  `json:"date"`
	ID        int    `json:"id"`
	Important bool   `json:"important"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
}
