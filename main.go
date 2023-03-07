package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type EbookInfoProduct struct {
    Title         string `json:"title"`
    Niche         string `json:"niche"`
    Organization string `json:"organizationId"`
    Product       string `json:"productId"`
    Format        string `json:"format"`
}


func createEbook(c *gin.Context) {
    title := c.PostForm("title")
    niche := c.PostForm("niche")
    organizationID := c.PostForm("organizationId")
    productID := c.PostForm("productId")

    ebook := EbookInfoProduct{
        Title:         title,
        Niche: niche,
        Organization: organizationID,
        Product:       productID,
        Format:        "ebook",
    }

    c.JSON(http.StatusOK, gin.H{"data": ebook})
}




func main() {

	router := gin.Default()


    // root path
	router.GET("/", func(c *gin.Context) {
        c.String(200, "The mentan blackbox is running...")
    })

    // generate ebook route
    router.POST("/api/v1/blackbox/ebook", createEbook)

	router.Run("localhost:8000")
}
