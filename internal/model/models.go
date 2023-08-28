package model

import "time"

type Segment struct {
	Slug       string     `json:"slug"`
	DeleteTime *time.Time `json:"delete_time"`
}

type UserSegments struct {
	UserID   uint64    `json:"user_id"`
	Segments []Segment `json:"segments"`
}
