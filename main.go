package main

import (
	"blogServer/controller"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main() {
	// Connect to database
	err := controller.ConnectDatabase()
	controller.CheckErr(err)

	r := gin.Default()

	// API router
	router := r.Group("/api")
	{
		router.GET("author", controller.GetAuthor)
		router.GET("author/:id", controller.GetAuthorById)
		router.POST("author", controller.AddAuthor)
		router.PUT("author/:id", controller.UpdateAuthor)
		router.DELETE("author/:id", controller.DeleteAuthor)

		router.GET("posts", controller.GetPost)
		router.GET("posts/:id", controller.GetPostById)
		router.POST("posts", controller.AddPosts)
		router.PUT("posts/:id", controller.UpdatePost)
		router.DELETE("posts/:id", controller.DeletePost)

		router.GET("comment", controller.GetComment)
		router.GET("comment/:id", controller.GetCommentById)
		router.POST("comment", controller.AddComment)
		router.PUT("comment/:id", controller.UpdateComment)
		router.DELETE("comment/:id", controller.DeleteComment)

		router.GET("tag", controller.GetTag)
		router.GET("tag/:id", controller.GetTagById)
		router.POST("tag", controller.AddTag)
		router.PUT("tag/:id", controller.UpdateTag)
		router.DELETE("tag/:id", controller.DeleteTag)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()

}






