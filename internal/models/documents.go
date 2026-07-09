package models

import "time"

type Document struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	Rubrics     []string  `json:"rubrics"`
	CreatedDate time.Time `json:"created_date"`
}

type DocumentWithoutID struct {
	Text        string    `json:"text"`
	Rubrics     []string  `json:"rubrics"`
	CreatedDate time.Time `json:"created_date"`
}
