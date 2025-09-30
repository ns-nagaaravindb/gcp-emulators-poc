package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	bucketName   = "test-bucket"
	gcsProjectID = "test-project"
)

func main() {
	// Set up context
	ctx := context.Background()

	// Create storage client pointing to the fake GCS server
	client, err := storage.NewClient(ctx, option.WithEndpoint("http://localhost:4443/storage/v1/"))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer client.Close()

	// Create bucket if it doesn't exist
	if err := createBucketIfNotExists(ctx, client); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	// Write a test file to GCS
	fileName := "test-file.txt"
	content := "Hello, this is a test file for GCS emulator!"

	if err := writeFile(ctx, client, fileName, content); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("Successfully wrote file '%s' to bucket '%s'\n", fileName, bucketName)

	// List files in the bucket
	fmt.Println("\nListing files in bucket:")
	if err := listFiles(ctx, client); err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	// Read the file back
	fmt.Printf("\nReading file '%s' content:\n", fileName)
	if err := readFile(ctx, client, fileName); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
}

func createBucketIfNotExists(ctx context.Context, client *storage.Client) error {
	bucket := client.Bucket(bucketName)

	// Check if bucket exists
	_, err := bucket.Attrs(ctx)
	if err == nil {
		fmt.Printf("Bucket '%s' already exists\n", bucketName)
		return nil
	}

	// Create bucket
	if err := bucket.Create(ctx, gcsProjectID, nil); err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	fmt.Printf("Created bucket '%s'\n", bucketName)
	return nil
}

func writeFile(ctx context.Context, client *storage.Client, fileName, content string) error {
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)

	// Create a writer
	w := obj.NewWriter(ctx)
	defer w.Close()

	// Write content
	if _, err := io.Copy(w, strings.NewReader(content)); err != nil {
		return fmt.Errorf("failed to write content: %v", err)
	}

	return nil
}

func listFiles(ctx context.Context, client *storage.Client) error {
	bucket := client.Bucket(bucketName)

	it := bucket.Objects(ctx, nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate objects: %v", err)
		}

		fmt.Printf("- %s (size: %d bytes, created: %v)\n",
			obj.Name, obj.Size, obj.Created.Format(time.RFC3339))
	}

	return nil
}

func readFile(ctx context.Context, client *storage.Client, fileName string) error {
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(fileName)

	// Create a reader
	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reader: %v", err)
	}
	defer r.Close()

	// Read content
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read content: %v", err)
	}

	fmt.Printf("Content: %s\n", string(content))
	return nil
}
