package firebase

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
)


type Introduction struct {
    ProductID    string
    Introduction string
}

var (
    firestoreOnce sync.Once
    firestoreClient *firestore.Client
    introductionsCollection string
)

func init() {
    introductionsCollection = os.Getenv("INTRODUCTIONS_PRODUCTS_COLLECTION")
}


func SaveIntroduction(intro Introduction) (bool, error) {
    ctx := context.Background()
    projectId := os.Getenv("FIREBASE_APP_ID")

    firestoreClient, err := firestore.NewClient(ctx, projectId)
    if err != nil {
        log.Fatalf("Failed to create a Firestore client: %v", err)
    }

    _, _, err = firestoreClient.Collection(introductionsCollection).Add(ctx, map[string]interface{}{
        "productId":    intro.ProductID,
        "introduction": intro.Introduction,
    })

    if err != nil {
        return false, err
    }

    return true, nil
}

