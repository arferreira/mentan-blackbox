package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)



func ConnectFirestore(ctx context.Context) (*firestore.Client, error) {

        var err error
        projectId := os.Getenv("FIREBASE_APP_ID")
        firestoreClient, err = firestore.NewClient(ctx, projectId)
        if err != nil {
            log.Fatalf("Failed to create a Firestore client: %v", err)
        }

    return firestoreClient, nil
}
