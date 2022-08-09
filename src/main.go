package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"database/sql"
	"example/irisuinwl/web-service-gin/db"
	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 60*time.Second)
		defer cancel()
		rows, err := db.QueryContext(ctx, `SELECT id, title, artist, price FROM albums;`)
		if err != nil {
			log.Panicf("query all users: %v", err)
		}
		defer rows.Close()
		var albums []*Album
		for rows.Next() {
			var (
				albumID, albumTitle, albumArtist string
				albumPrice                       float64
			)
			if err := rows.Scan(&albumID, &albumTitle, &albumArtist, &albumPrice); err != nil {
				log.Panicf("scan the album: %v", err)
			}
			albums = append(albums, &Album{
				ID:     albumID,
				Title:  albumTitle,
				Artist: albumArtist,
				Price:  albumPrice,
			})
		}
		if err := rows.Close(); err != nil {
			log.Panicf("rows close: %v", err)
		}

		if err := rows.Err(); err != nil {
			log.Panicf("scan users: %v", err)
		}

		c.IndentedJSON(http.StatusOK, albums)
	}

}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 60*time.Second)
		defer cancel()

		var newAlbum Album
		// Call BindJSON to bind the received JSON to newAlbum.
		if err := c.BindJSON(&newAlbum); err != nil {
			log.Panicf("bind json failed: %v", err)
			return
		}
		log.Printf("newAlbum: %v", newAlbum)
		// Add the new album to db.
		query := `INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)
		RETURNING id, title, artist, price;
		`
		row := db.QueryRowContext(
			ctx, query,
			newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price,
		)
		var createdAlbum Album
		err := row.Scan(&createdAlbum.ID, &createdAlbum.Title, &createdAlbum.Artist, &createdAlbum.Price)
		if err != nil {
			log.Panicf("query all users: %v", err)
		}

		c.IndentedJSON(http.StatusCreated, newAlbum)
	}

}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(c, 60*time.Second)
		defer cancel()
		row := db.QueryRowContext(ctx, `SELECT id, title, artist, price FROM albums WHERE id = $1;`, id)
		var (
			albumID, albumTitle, albumArtist string
			albumPrice                       float64
		)
		err := row.Scan(&albumID, &albumTitle, &albumArtist, &albumPrice)
		if err != nil {
			log.Panicf("scan the album: %v", err)
		}
		album := Album{
			ID:     albumID,
			Title:  albumTitle,
			Artist: albumArtist,
			Price:  albumPrice,
		}
		c.IndentedJSON(http.StatusOK, album)
	}
}

func main() {
	router := gin.Default()
	handler := db.NewHandler()
	db := handler.GetDB()

	router.GET("/albums", getAlbums(db))
	router.GET("/albums/:id", getAlbumByID(db))
	router.POST("/albums", postAlbums(db))

	router.Run("0.0.0.0:8080")
}
