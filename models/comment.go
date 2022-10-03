package models

type Comment struct {
	Id       int    `json:"id"`
	Comment  string `json:"comment"`
	AuthorID string `json:"author_id"`
	PostId   string `json:"post_id"`
}
