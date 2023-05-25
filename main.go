package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

type Track struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title    string `json:"title"`
		CoverURL string `json:"cover_big"`
	} `json:"album"`
	Duration   int    `json:"duration"`
	PreviewURL string `json:"preview"`
}

type TracksResponse struct {
	Data []Track `json:"data"`
}

type TrackWithLikeCount struct {
	Track     Track
	LikeCount int
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("view/*.html")
	router.Static("/css", "./view/css")

	// Open a connection to the SQL database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/DeezDB")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/", func(c *gin.Context) {
		limit := 50 // Set the desired limit to retrieve more tracks, e.g., 50
		indexStr := c.Query("index")
		index, _ := strconv.Atoi(indexStr)

		// Construct the URL with pagination parameters
		url := fmt.Sprintf("https://api.deezer.com/chart/0/tracks?limit=%d&index=%d", limit, index)

		// Retrieve data from Deezer API
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		var tracksResponse TracksResponse
		err = json.NewDecoder(response.Body).Decode(&tracksResponse)
		if err != nil {
			log.Fatal(err)
		}

		// Get like counts for each track from the database
		tracksWithLikeCount := make([]TrackWithLikeCount, 0)
		for _, track := range tracksResponse.Data {
			likeCount := getLikeCount(track.ID, db)
			tracksWithLikeCount = append(tracksWithLikeCount, TrackWithLikeCount{
				Track:     track,
				LikeCount: likeCount,
			})
		}

		// Render data in an HTML table
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Tracks": tracksWithLikeCount,
		})
	})

	router.POST("/like", func(c *gin.Context) {
		var requestData struct {
			TrackID int `json:"trackID"`
		}

		err := c.ShouldBindJSON(&requestData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Increment the like count in the database
		err = incrementLikeCount(requestData.TrackID, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment like count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	log.Fatal(router.Run(":8080"))
}

func incrementLikeCount(trackID int, db *sql.DB) error {
	// Check if the track exists in the database
	query := "SELECT like_counts FROM tracks WHERE track_id = ?"
	var currentLikeCount int
	err := db.QueryRow(query, trackID).Scan(&currentLikeCount)
	if err != nil {
		if err == sql.ErrNoRows {
			// Track does not exist, insert a new row
			_, err = db.Exec("INSERT INTO tracks (track_id, like_counts) VALUES (?, 1)", trackID)
			return err
		}
		return err
	}

	// Increment the like count
	_, err = db.Exec("UPDATE tracks SET like_counts = ? WHERE track_id = ?", currentLikeCount+1, trackID)
	return err
}

func getLikeCount(trackID int, db *sql.DB) int {
	query := "SELECT like_counts FROM tracks WHERE track_id = ?"
	var likeCount int
	err := db.QueryRow(query, trackID).Scan(&likeCount)
	if err != nil {
		if err == sql.ErrNoRows {
			// Track does not exist, return 0 like count
			return 0
		}
		log.Println("Failed to retrieve like count:", err)
		return 0
	}
	return likeCount
}
