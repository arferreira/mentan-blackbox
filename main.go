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
    chaptersPrompt := fmt.Sprintf("Create 10 chapters for an ebook for me about %s with this %s niche using correct grammar and engaged words. A title consists of a litle and the number of the chapter, separeted by colon", ebook.Title, ebook.Niche)
    chapters := openai.OpenAiPrompt{
        Prompt:      chaptersPrompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(chapters.Prompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate chapter content: %v", err)
    }
    return response, nil
}

// get introductions
func getIntroduction(ebook EbookInfoProduct) (string, error) {
    introductionPrompt := fmt.Sprintf("Create a short introduction for an ebook for me about %s with this %s niche using correct grammar and engaged words.", ebook.Title, ebook.Niche)

    introduction := openai.OpenAiPrompt{
        Prompt:      introductionPrompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(introduction.Prompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate introduction content: %v", err)
    }
    return response, nil
}

func getChaptersContent(ebook EbookInfoProduct, chapter string) (string, error) {
    chapterContentPrompt := fmt.Sprintf("Teach a student about the below topic and subtopic and by writing multiple paragraphs. the topix is: %s and the subtopic is: %s and the name of the chapter is: %s", ebook.Title, ebook.Niche, chapter)

    chapterContent := openai.OpenAiPrompt{
        Prompt:      chapterContentPrompt,
        Temperature: 1.0,
        MaxTokens:   10,
    }

    response, err := openai.SecondLayer(chapterContent.Prompt)
    if err != nil {
        return "", fmt.Errorf("failed to generate chapterContent content: %v", err)
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

    introduction, err := getIntroduction(ebook)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate introduction content"})
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


   // create array of chapters
    chaptersArray := make([]string, 0, len(chapters))
    for _, chapter := range chapters {
        chapterContent, err := getChaptersContent(ebook, string(chapter))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate chapter content"})
            fmt.Fprintf(os.Stderr, "%v\n", err)
            return
        }
        chaptersArray = append(chaptersArray, chapterContent)

        // save chapter content on firebase
        _, err = saveProductChapter(ebook.Product, ebook.Organization, chapter, chapterContent)
    }

    // Save all chapters in collection products
    err = updateProductChapters(ebook.Product, ebook.Organization, chaptersArray)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chapter contents"})
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }


    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "introduction": introduction,
            "chapters":     chapterStrings,
        },
    })


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
