package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("view/*.html")
	router.Static("/css", "./view/css")

	router.GET("/", handleRequest)

	log.Fatal(router.Run(":8080"))
}

func handleRequest(c *gin.Context) {
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

	// Render data in an HTML table
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Tracks": tracksResponse.Data,
	})
}
