package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Ctag represents a single universal-ctags JSON output record.
type Ctag struct {
	Type     string `bson:"_type"`
	Name     string `bson:"name"`
	Path     string `bson:"path"`
	Pattern  string `bson:"pattern"`
	Language string `bson:"language"`
	Line     int    `bson:"line"`
	Kind     string `bson:"kind"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <connection-uri> <json-file>\n", os.Args[0])
		os.Exit(1)
	}

	uri := os.Args[1]
	jsonFile := os.Args[2]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	db := client.Database("ctags")
	col := db.Collection("ctags")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "path", Value: 1},
			{Key: "line", Value: 1},
			{Key: "name", Value: 1},
			{Key: "kind", Value: 1},
			{Key: "pattern", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, indexModel); err != nil {
		log.Fatalf("failed to create index: %v", err)
	}

	jobs := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go insertTags(ctx, col, jobs, &wg)

	file, err := os.Open(jsonFile)
	if err != nil {
		log.Fatalf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read line: %v", err)
		}
		jobs <- line
	}
	close(jobs)
	wg.Wait()
}

func insertTags(ctx context.Context, col *mongo.Collection, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		var doc Ctag
		if err := bson.UnmarshalExtJSON([]byte(j), true, &doc); err != nil {
			log.Printf("failed to decode line: %v", err)
			continue
		}
		filter := bson.M{
			"path":    doc.Path,
			"line":    doc.Line,
			"name":    doc.Name,
			"kind":    doc.Kind,
			"pattern": doc.Pattern,
		}
		_, err := col.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
		if err != nil {
			log.Printf("failed to insert document: %v", err)
		}
	}
}
