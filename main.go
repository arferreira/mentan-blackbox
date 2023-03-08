package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arferreira/mentan-blackbox/openai"
	"github.com/gin-contrib/cors"
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


type ebookData struct {
    Title           string `json:"title"`
    Niche           string `json:"niche"`
    OrganizationID  string `json:"organizationId"`
    ProductID       string `json:"productId"`
}


func createEbook(c *gin.Context) {
    fmt.Println(c.Request.Body)

    var data ebookData

    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    title := data.Title
    niche := data.Niche
    organizationID := data.OrganizationID
    productID := data.ProductID

    log.Printf("Title: %s", title)
    log.Printf("Niche: %s", niche)
    log.Printf("Organization ID: %s", organizationID)
    log.Printf("Product ID: %s", productID)




    ebook := EbookInfoProduct{
        Title:         title,
        Niche: niche,
        Organization: organizationID,
        Product:       productID,
        Format:        "ebook",
    }


    prompt := fmt.Sprintf("Create a ebook for me about %s with this %s niche using correct grammar and engaged content. I want it with introduction, 10 chapters and conclusion.", ebook.Title, ebook.Niche)

    secondLayer := openai.OpenAiPrompt{
        Prompt:      prompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(secondLayer.Prompt)
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
    router.Use(gin.Logger())

    // Configure CORS middlewares
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"}
    router.Use(cors.New(config))


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
