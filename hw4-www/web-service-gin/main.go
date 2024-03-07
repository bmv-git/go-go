package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
//   - gin.Context is the most important part of Gin.
//     It carries request details, validates and serializes JSON, and more.
//     *) Despite the similar name, this is different from Go’s built-in context package.
//   - Context.IndentedJSON to serialize the struct into JSON and add it to the response.
//     The function’s first argument is the HTTP status code to send to the client.
//     Here, the StatusOK constant from the net/http package is being passed
//     to indicate 200 OK.
//     *) Note that Context.IndentedJSON can be replaced with a call to Context.JSON
//     to send more compact JSON. In practice, the indented form is much easier to work with
//     when debugging and the size difference is usually small.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Calling BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

//   - Initializing a Gin router using Default.
//   - Using the GET function to associate the GET HTTP method and /albums path
//     with a handler function.
//     *) Note that you’re passing the name of the getAlbums function.
//     This is different from passing the result of the function,
//     which you would do by passing getAlbums() (note the parenthesis).
//   - Using the Run function to attach the router to an http.Server and start the server.
func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	err := router.Run("localhost:8000")
	if err != nil {
		return
	}
}
