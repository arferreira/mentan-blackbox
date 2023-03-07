package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Fake data
type user struct {
    Id        int     `json:"id"`
    Username  string  `json:"username"`
}

// feed the user structure
var users = []user{
    {Id: 546, Username: "John"},
    {Id: 894, Username: "Mary"},
    {Id: 326, Username: "Jane"},
}


// GET /users rule
func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}




func main() {



	router := gin.Default()

	// TODO: define another routes
	router.GET("/", func(c *gin.Context) {
        c.String(200, "The mentan blackbox is running...")
    })
	router.GET("/users", getUsers)

	router.Run("localhost:8000")
}
