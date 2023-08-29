package model

import (
	"fmt"
	"time"
)

type OpType int

const (
	AddOp = OpType(iota)
	DeleteOp
)

func (o OpType) String() string {
	if o == AddOp {
		return "add"
	}
	return "delete"
}

type UserSegment struct {
	UserID     uint64     `json:"user_id"`
	Slug       string     `json:"slug"`
	DeleteTime *time.Time `json:"delete_time"`
}

type UserLog struct {
	UserID      uint64    `json:"user_id"`
	Slug        string    `json:"slug"`
	Operation   string    `json:"operation"`
	RequestTime time.Time `json:"request_time"`
}

func (u UserLog) String() string {
	return fmt.Sprintf("%d;%s;%s;%v", u.UserID, u.Slug, u.Operation, u.RequestTime)
}
