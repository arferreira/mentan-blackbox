package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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


// getChapters generates ebook chapters using OpenAI API
func getChapters(ebook EbookInfoProduct) (string, error){
    chaptersPrompt := fmt.Sprintf("Create 10 chapters for an ebook for me about %s with this %s niche using correct grammar and engaged words. Give me separated by comon", ebook.Title, ebook.Niche)
    chapters := openai.OpenAiPrompt{
        Prompt:      chaptersPrompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(chapters.Prompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate ebook content: %v", err)
    }
    return response, nil
}

func getDescription(ebook EbookInfoProduct) (string, error) {
    descriptionPrompt := fmt.Sprintf("Create a description for an ebook for me about %s with this %s niche using correct grammar and engaged words.", ebook.Title, ebook.Niche)

    description := openai.OpenAiPrompt{
        Prompt:      descriptionPrompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(description.Prompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate ebook content: %v", err)
    }
    return response, nil
}

func createEbook(c *gin.Context) {
    startTime := time.Now()

    var data ebookData
    log.Printf("data: %v", data)

    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        fmt.Println("Err: ", err)
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
        Niche:         niche,
        Organization: organizationID,
        Product:       productID,
        Format:        "ebook",
    }

    description, err := getDescription(ebook)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ebook content"})
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }

    chapters, err := getChapters(ebook)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ebook content"})
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }
    fmt.Println("Chapters generated", chapters)
    c.JSON(http.StatusOK, gin.H{"data": description})

    elapsedSeconds := time.Since(startTime).Seconds()
    log.Printf("Elapsed time for the first layer: %f seconds", elapsedSeconds)
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
