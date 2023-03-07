package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arferreira/mentan-blackbox/openai"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

    if len(title) == 0 || len(niche) == 0 || len(organizationID) == 0 || len(productID) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
        return
    }

    ebook := EbookInfoProduct{
        Title:         title,
        Niche: niche,
        Organization: organizationID,
        Product:       productID,
        Format:        "ebook",
    }


    prompt := fmt.Sprintf("Create a ebook for me about %s using this %s", ebook.Niche, ebook.Title)

    secondLayer := openai.OpenAiPrompt{
        Prompt:      prompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SendPrompt(secondLayer.Prompt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ebook content"})
        fmt.Fprintf(os.Stderr, "failed to generate ebook content: %v\n", err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": response})
}




func main() {

    // load .env file
    err := godotenv.Load()
    if err != nil {
    log.Fatal("Error loading .env file")
    }


	router := gin.Default()

    // root path
	router.GET("/", func(c *gin.Context) {
        c.String(200, "The mentan blackbox is running...")
    })

    // generate ebook route
    router.POST("/api/v1/blackbox/ebook", createEbook)

	if err := router.Run(":8000"); err != nil {
        fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
        os.Exit(1)
    }
}
