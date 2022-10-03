package models


type Post struct {
	Id        int    `json:"id"`
	Title string `json:"title"`
	Description  string `json:"description"`
	AuthorID     string `json:"author_id"`
	TagList     []string `json:"tags"`
}

type PostView struct {
	Id        int    `json:"id"`
	Title string `json:"title"`
	Description  string `json:"description"`
	AuthorID     string `json:"author_id"`
}