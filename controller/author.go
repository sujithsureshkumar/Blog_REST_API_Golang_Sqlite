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

//getting all authors

func AuthorsGet() ([]models.Author, error) {

	rows, err := DB.Query("SELECT id, first_name, last_name, email from authors")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	author := make([]models.Author, 0)

	for rows.Next() {
		singleAuthor := models.Author{}
		err = rows.Scan(&singleAuthor.Id, &singleAuthor.FirstName, &singleAuthor.LastName, &singleAuthor.Email)

		if err != nil {
			return nil, err
		}

		author = append(author, singleAuthor)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return author, err
}

func GetAuthor(c *gin.Context) {

	authors, err := AuthorsGet()

	CheckErr(err)

	if authors == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": authors})
	}
}

//get author by Id
func AuthorGetById(id string) (models.Author, error) {

	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email from authors WHERE id = ?")

	if err != nil {
		return models.Author{}, err
	}

	author := models.Author{}

	sqlErr := stmt.QueryRow(id).Scan(&author.Id, &author.FirstName, &author.LastName, &author.Email)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.Author{}, nil
		}
		return models.Author{}, sqlErr
	}
	return author, nil
}

func GetAuthorById(c *gin.Context) {

	// grab the Id of the record want to retrieve
	id := c.Param("id")

	author, err := AuthorGetById(id)

	CheckErr(err)
	// if the name is blank we can assume nothing is found
	if author.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": author})
	}
}

//Add a new Author

func AddAuthor(c *gin.Context) {
	var json models.Author

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := AuthorAdd(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}


func AuthorAdd(newAuthor models.Author) (bool, error) {
   
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO authors (first_name, last_name, email) VALUES (?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newAuthor.FirstName, newAuthor.LastName, newAuthor.Email)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}


//update Author by id
func AuthorUpdate(newAuthor models.Author, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE authors SET first_name = ?, last_name = ?, email = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newAuthor.FirstName, newAuthor.LastName, newAuthor.Email, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateAuthor(c *gin.Context) {

	var json models.Author

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorId, err := strconv.Atoi(c.Param("id"))

	fmt.Printf("Updating id %d", authorId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := AuthorUpdate(json, authorId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

//Delete Author by Id
func AuthorDelete(authorId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("DELETE from authors where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(authorId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteAuthor(c *gin.Context) {

	authorId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := AuthorDelete(authorId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}