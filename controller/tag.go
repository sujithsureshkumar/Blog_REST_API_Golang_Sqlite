package controller

import (
	"blogServer/models"
	"fmt"
	"net/http"
	"strconv"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

//Getting all Tags
func TagsGet() ([]models.Tag, error) {

	rows, err := DB.Query("SELECT id, tag, author_id,post_id from tags")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tag := make([]models.Tag, 0)

	for rows.Next() {
		singleTag := models.Tag{}
		err = rows.Scan(&singleTag.Id, &singleTag.Tag, &singleTag.AuthorID, &singleTag.PostId)

		if err != nil {
			return nil, err
		}

		tag = append(tag, singleTag)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return tag, err
}

func GetTag(c *gin.Context) {

	tags, err := TagsGet()

	CheckErr(err)

	if tags == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": tags})
	}
}

//get Tag by Id
func TagGetById(id string) (models.Tag, error) {

	stmt, err := DB.Prepare("SELECT id, tag, post_id, author_id from tags WHERE id = ?")

	if err != nil {
		return models.Tag{}, err
	}

	tag := models.Tag{}

	sqlErr := stmt.QueryRow(id).Scan(&tag.Id, &tag.Tag, &tag.PostId, &tag.AuthorID)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Tag{}, nil
		}
		return models.Tag{}, sqlErr
	}
	return tag, nil
}

func GetTagById(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	tag, err := TagGetById(id)

	CheckErr(err)
	// if the name is blank we can assume nothing is found
	if tag.Tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": tag})
	}
}

func AddTag(c *gin.Context) {
	var json models.Tag

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := TagAdd(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func TagAdd(newTag models.Tag) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO tags (tag,author_id,post_id) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newTag.Tag, newTag.AuthorID, newTag.PostId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

// update Tag by id
func TagUpdate(newTag models.Tag, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE tags SET tag = ?  WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newTag.Tag, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateTag(c *gin.Context) {

	var json models.Tag

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", tagId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := TagUpdate(json, tagId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

//Delete Tag by id
func TagDelete(tagId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from tags where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(tagId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteTag(c *gin.Context) {

	tagId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := TagDelete(tagId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
