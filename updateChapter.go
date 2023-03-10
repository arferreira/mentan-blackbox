package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

func updateChapters(ctx context.Context, productID string, chapters []string) error {
    // Get reference to the Firestore client
    client, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_APP_ID"))
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }
    defer client.Close()

    // Reference to the document that needs to be updated
    docRef := client.Collection("products").Doc(productID)

    // Update the "chapters" field with new chapters
    _, err = docRef.Update(ctx, []firestore.Update{
        {Path: "chapters", Value: chapters},
    })
    if err != nil {
        return fmt.Errorf("failed to update chapters for product ID %s: %v", productID, err)
    }

    log.Printf("Successfully updated chapters for product ID %s", productID)
    return nil
}


func saveProductChapter(chapterNum int, productName string, chapterContent string) error {
    // get access to Firestore db
    ctx := context.Background()
    projectID := os.Getenv("FIREBASE_APP_ID")
    firestoreClient, err := firestore.NewClient(ctx, projectID)
    if err != nil {
        return fmt.Errorf("failed to create Firestore client: %v", err)
    }
    defer firestoreClient.Close()

    // build reference to the products collection and the specific product document by name
    productDocRef := firestoreClient.Collection("products").Doc(productName)

    // update the chapters field of the product document to include the new chapter content
    _, err = productDocRef.Update(ctx, []firestore.Update{
        {Path: fmt.Sprintf("chapters.chapter%d", chapterNum), Value: chapterContent},
    })
    if err != nil {
        return fmt.Errorf("failed to update chapters on product document %s: %v", productName, err)
    }

    return nil
}
