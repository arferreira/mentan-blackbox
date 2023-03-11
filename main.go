package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/arferreira/mentan-blackbox/openai"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	firebase "firebase.google.com/go"
)

type EbookInfoProduct struct {
	Title        string `json:"title"`
	Niche        string `json:"niche"`
	Organization string `json:"organizationId"`
	Product      string `json:"productId"`
	Format       string `json:"format"`
}

type ebookData struct {
	Title          string `json:"title"`
	Niche          string `json:"niche"`
	OrganizationID string `json:"organizationId"`
	ProductID      string `json:"productId"`
}

// define constant values
const (
	numWorkers = 5
	bufferSize = 20
)

// generatePrompt generates a prompt using the given format and parameters,
// and passes it to the OpenAI API to obtain a response.
func generatePrompt(format string, args ...interface{}) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(format, args...))

	prompt := openai.OpenAiPrompt{
		Prompt:      sb.String(),
		Temperature: 1.0,
		MaxTokens:   5,
	}

	response, err := openai.SecondLayer(prompt.Prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt content: %v", err)
	}
	return response, nil
}

// getChapters generates ebook chapters using OpenAI API with go routine
func getChapters(ctx context.Context, ebook EbookInfoProduct, ch chan<- string, numConcurrentJobs int) {
	// create buffered channel for limiting a number of concurrent jobs
	semaphore := make(chan bool, numConcurrentJobs)

	for i := 0; i < 10; i++ {

		semaphore <- true // write to channel - blocking if channel is full

		go func() {
			defer func() { <-semaphore }() // read from channel

			chapterPrompt := fmt.Sprintf("Create chapter %d for '%s' focused in this '%s' niche.", i+1, ebook.Title, ebook.Niche)
			fmt.Println("Chapter number: ", i+1)
			chapterContent, err := openai.SecondLayer(chapterPrompt)
			if err != nil {
				log.Printf("[ERROR] failed to generate chapter content for prompt '%s': %v", chapterPrompt, err)
				return
			}
			ch <- chapterContent
		}()
	}

	// wait until all running goroutines finish
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}
}

// getIntroduction generates a short introduction for an ebook using the given product info.
func getIntroduction(ebook EbookInfoProduct) (string, error) {
	format := "Create a short introduction for an ebook for me about %s with this %s niche using correct grammar and engaged words."
	return generatePrompt(format, ebook.Title, ebook.Niche)
}

// getChaptersContent generates the content for a chapter in an ebook using the given product info and chapter name.
func getChaptersContent(ebook EbookInfoProduct, chapter string) (string, error) {
	format := "Teach a student about the below topic and subtopic and by writing multiple paragraphs. The topic is: %s and the subtopic is: %s and the name of the chapter is: %s"
	return generatePrompt(format, ebook.Title, ebook.Niche, chapter)
}

func createEbook(ctx context.Context, c *gin.Context) {

	startTime := time.Now()

	var data ebookData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title := data.Title
	niche := data.Niche
	organizationID := data.OrganizationID
	productID := data.ProductID

	ebook := EbookInfoProduct{
		Title:        title,
		Niche:        niche,
		Organization: organizationID,
		Product:      productID,
		Format:       "ebook",
	}

	// Generate introduction content.
	introduction, err := getIntroduction(ebook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate introduction content"})
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Create a buffered channel for receiving chapter names.
	chapterChan := make(chan string, bufferSize)

	// Use a context with a deadline to cancel long-running operations.
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Use a wait group to wait for all workers to finish.
	wg := new(sync.WaitGroup)

	// Start worker goroutines to generate chapter content.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getChapters(ctx, ebook, chapterChan, bufferSize)
		}()
	}

	// Collect the chapter names from the channel.
	chapters := []string{}
	for chapter := range chapterChan {
		chapters = append(chapters, chapter)
	}
	close(chapterChan)

	// Check if any chapters were received.
	if len(chapters) == 0 {
		err = fmt.Errorf("failed to receive any chapters")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate chapter content for each chapter name.
	chaptersArray := make([]string, 0, len(chapters))
	for _, chapter := range chapters {
		chapterContent, err := getChaptersContent(ebook, chapter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate chapter content"})
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		chaptersArray = append(chaptersArray, chapterContent)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the ebook"})
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get firestore client"})
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get firestore client"})
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	defer client.Close()
	collection := client.Collection("products")
	docRef := collection.Doc(productID)

	dataToSave := gin.H{
		"title":        title,
		"niche":        niche,
		"organization": organizationID,
		"product":      productID,
		"format":       "ebook",
		"introduction": introduction,
		"chapters":     chaptersArray,
	}

	_, err = docRef.Set(ctx, dataToSave)
	if err != nil {
		// get the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ebook"})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"introduction": introduction,
			"chapters":     chaptersArray,
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

	createEbook := func(c *gin.Context) {
		ctx := context.Background()
		createEbook(ctx, c)
	}

	// generate ebook route
	router.POST("/api/v1/blackbox/ebook", createEbook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
		os.Exit(1)
	}
}
