package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		DBName:               os.Getenv("DBNAME"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DBHOST"), os.Getenv("DBPORT")),
		AllowNativePasswords: true,
	}

	maxRetries := 5
	retryInterval := 5 * time.Second

	var err error

	for retry := 0; retry < maxRetries; retry++ {
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			fmt.Printf("Failed to connect to the database: %v\n", err)
			fmt.Printf("Retrying in %s...\n", retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		if err := db.Ping(); err != nil {
			fmt.Printf("Failed to ping the database: %v\n", err)
			fmt.Printf("Retrying in %s...\n", retryInterval)
			db.Close()
			time.Sleep(retryInterval)
			continue
		}

		defer db.Close()
		break
	}

	if err != nil {
		fmt.Errorf("Failed to connect after %d retries: %v\n", maxRetries, err)
		return
	}

	log.Println("Connected to MariaDB.")

	router := gin.Default()
	router.GET("/albums", getAllAlbums)
	router.GET("/album", getAlbumsByArtist)
	router.GET("/album/:id", albumByID)
	router.POST("/album", postAlbum)

	router.Run("0.0.0.0:8080")
}

func getAllAlbums(c *gin.Context) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("allAlbums %v", err)})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("allAlbums %v", err)})
			return
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("allAlbums %v", err)})
		return
	}

	if len(albums) == 0 {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumsByArtist(c *gin.Context) {
	artist := c.Query("artist")

	if artist == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter \"artist\" is required."})
		return
	}

	albums, err := privateAlbumsByArtist(artist)
	if albums == nil || err != nil {
		message := fmt.Sprintf("album of artist = %s not found.", artist)
		c.JSON(http.StatusNotFound, gin.H{"message": message})
		return
	}

	if albums != nil {
		c.IndentedJSON(http.StatusOK, albums)
	}
}

func privateAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM albums WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("privateAlbumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("privateAlbumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("privateAlbumsByArtist %q, %v", name, err)
	}
	return albums, nil
}

func postAlbum(c *gin.Context) {
	var alb Album
	if err := c.BindJSON(&alb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("postAlbum %v", err)})
		return
	}

	result, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("postAlbum %v", err)})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("postAlbum %v", err)})
		return
	}

	row := db.QueryRow("SELECT * FROM albums WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("postAlbum %v", err)})
		return
	}
	c.JSON(http.StatusOK, alb)
}

func albumByID(c *gin.Context) {
	var alb Album
	id := c.Param("id")

	row := db.QueryRow("SELECT * FROM albums WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("albumById %v", err)})
		return
	}
	c.JSON(http.StatusOK, alb)
}
