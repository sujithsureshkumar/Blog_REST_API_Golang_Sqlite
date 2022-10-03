package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	   "blogServer/models"
	   "database/sql"
	   "net/http"
	   "strconv"
	   "fmt"
	

)

//Getting all Posts
func PostsGet() ([]models.PostView, error) {

	rows, err := DB.Query("SELECT id, title, description, author_id from posts")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	post := make([]models.PostView, 0)

	for rows.Next() {
		singlePost := models.PostView{}
		err = rows.Scan(&singlePost.Id, &singlePost.Title, &singlePost.Description, &singlePost.AuthorID)

		if err != nil {
			return nil, err
		}

		post = append(post, singlePost)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return post, err
}

func GetPost(c *gin.Context) {

	posts, err := PostsGet()

	CheckErr(err)

	if posts == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": posts})
	}
}

//get post by Id
func PostGetById(id string) (models.PostView, error) {

	stmt, err := DB.Prepare("SELECT id, title, description,author_id from posts WHERE id = ?")

	if err != nil {
		return models.PostView{}, err
	}

	post := models.PostView{}

	sqlErr := stmt.QueryRow(id).Scan(&post.Id, &post.Title, &post.Description, &post.AuthorID)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.PostView{}, nil
		}
		return models.PostView{}, sqlErr
	}
	return post, nil
}

func GetPostById(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	post, err := PostGetById(id)

	CheckErr(err)
	// if the name is blank we can assume nothing is found
	if post.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}

//add a new post

func AddPosts(c *gin.Context) {
	var json models.Post

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := PostsAdd(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}


func PostsAdd(newPost models.Post) (bool, error) {
   
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO posts (title, description, author_id) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()
	var id int64
   var result sql.Result
	result, err = stmt.Exec(newPost.Title, newPost.Description, newPost.AuthorID)
	//fmt.Println(len(newPost.TagList))
	id,_ =result.LastInsertId()
	fmt.Println(id)

	if err != nil {
		return false, err
	}
	var stmt2 *sql.Stmt
	stmt2,_ =tx.Prepare("INSERT INTO tags (tag,author_id,post_id) VALUES (?, ?, ?)")
	defer stmt2.Close()
	for i :=0; i<len(newPost.TagList); i++{
		_, _ = stmt2.Exec(newPost.TagList[i], newPost.AuthorID, id)
	}

	tx.Commit()

	return true, nil
}

//update Post by id
func PostUpdate(newPost models.Post, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE posts SET title = ?,description = ?  WHERE Id = ? AND author_id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPost.Title,newPost.Description, id,newPost.AuthorID)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdatePost(c *gin.Context) {

	var json models.Post

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", postId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := PostUpdate(json, postId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

//Delete Post by id
func PostDelete(postId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from posts where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(postId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeletePost(c *gin.Context) {

	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := PostDelete(postId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}