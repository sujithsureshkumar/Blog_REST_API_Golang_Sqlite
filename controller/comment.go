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
func CommentsGet() ([]models.Comment, error) {

	rows, err := DB.Query("SELECT id, comment, author_id,post_id from comments")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comment := make([]models.Comment, 0)

	for rows.Next() {
		singleComment := models.Comment{}
		err = rows.Scan(&singleComment.Id, &singleComment.Comment, &singleComment.AuthorID, &singleComment.PostId)

		if err != nil {
			return nil, err
		}

		comment = append(comment, singleComment)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return comment, err
}

func GetComment(c *gin.Context) {

	comments, err := CommentsGet()

	CheckErr(err)

	if comments == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": comments})
	}
}

//get comment by Id
func CommentGetById(id string) (models.Comment, error) {

	stmt, err := DB.Prepare("SELECT id, comment,post_id,author_id from comments WHERE id = ?")

	if err != nil {
		return models.Comment{}, err
	}

	comment := models.Comment{}

	sqlErr := stmt.QueryRow(id).Scan(&comment.Id, &comment.Comment, &comment.PostId, &comment.AuthorID)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Comment{}, nil
		}
		return models.Comment{}, sqlErr
	}
	return comment, nil
}

func GetCommentById(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	comment, err := CommentGetById(id)

	CheckErr(err)
	// if the name is blank we can assume nothing is found
	if comment.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": comment})
	}
}

//Comment by id

func AddComment(c *gin.Context) {
	var json models.Comment

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := CommentAdd(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}


func CommentAdd(newComment models.Comment) (bool, error) {
   
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO comments (comment,author_id,post_id) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newComment.Comment, newComment.AuthorID, newComment.PostId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

//update Comment by id
func CommentUpdate(newComment models.Comment, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE comments SET comment = ?  WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newComment.Comment, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateComment(c *gin.Context) {

	var json models.Comment

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", commentId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := CommentUpdate(json, commentId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

//Delete Comment by id
func CommentDelete(commentId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from comments where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(commentId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteComment(c *gin.Context) {

	commentId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := CommentDelete(commentId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}