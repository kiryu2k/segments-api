package model

import "time"

const (
	AddOp = iota
	DeleteOp
)

type UserSegment struct {
	UserID     uint64     `json:"user_id"`
	Slug       string     `json:"slug"`
	DeleteTime *time.Time `json:"delete_time"`
}
