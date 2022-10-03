package models

type Tag struct {
	Id       int    `json:"id"`
	Tag      string `json:"tag"`
	PostId   string `json:"post_id"`
	AuthorID string `json:"author_id"`
}
