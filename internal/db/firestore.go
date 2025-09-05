package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreClient struct {
	client *firestore.Client
}

func NewFirestoreClient(ctx context.Context, projectId, databaseId string) (*FirestoreClient, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectId, databaseId)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}
	return &FirestoreClient{client: client}, nil
}

func (f *FirestoreClient) CheckDocumentExists(ctx context.Context, collection, documentId string) (bool, error) {
	ref := f.client.Collection(collection).Doc(documentId)
	doc, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}
	return doc.Exists(), nil
}

func (f *FirestoreClient) FetchDocument(ctx context.Context, collection, documentId string) (*firestore.DocumentSnapshot, error) {
	ref := f.client.Collection(collection).Doc(documentId)
	return ref.Get(ctx)
}

func (f *FirestoreClient) CreateOrOverwriteDocument(ctx context.Context, collection, documentId string, data interface{}) error {
	_, err := f.client.Collection(collection).Doc(documentId).Set(ctx, data)
	return err
}
