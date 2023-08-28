package model

import "time"

type UserSegment struct {
	UserID     uint64     `json:"user_id"`
	Slug       string     `json:"slug"`
	DeleteTime *time.Time `json:"delete_time"`
}

type ErrSegmentInfo struct {
	Slug    string `json:"slug"`
	Message string `json:"message"`
}
