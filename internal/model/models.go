package model

type UserSegments struct {
	UserID uint64   `json:"user_id"`
	Slugs  []string `json:"slugs"`
}
