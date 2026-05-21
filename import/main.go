package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CodeLine represents a single line of source code stored in the database.
type CodeLine struct {
	FilePath  string `bson:"path"`
	Line      string `bson:"line"`
	LineCount int    `bson:"line_count"`
}

func readFile(fileName string) ([]CodeLine, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var codeLines []CodeLine
	lineCount := 1
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		doc := CodeLine{LineCount: lineCount, Line: scanner.Text(), FilePath: fileName}
		codeLines = append(codeLines, doc)
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return codeLines, nil
}

func writeCode(ctx context.Context, col *mongo.Collection, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range jobs {
		codeLines, err := readFile(path)
		if err != nil {
			log.Printf("failed to read file %s: %v", path, err)
			continue
		}
		maxLineCount := 0
		for _, doc := range codeLines {
			if doc.LineCount > maxLineCount {
				maxLineCount = doc.LineCount
			}
			filter := bson.M{"path": doc.FilePath, "line_count": doc.LineCount}
			_, err := col.ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
			if err != nil {
				log.Printf("failed to insert document for %s: %v", path, err)
			}
		}
		_, err = col.DeleteMany(ctx, bson.M{
			"path":       path,
			"line_count": bson.M{"$gt": maxLineCount},
		})
		if err != nil {
			log.Printf("failed to trim stale lines for %s: %v", path, err)
		}
	}
}

// IsText returns true if the content appears to be text-based.
func IsText(content []byte) bool {
	contentType := http.DetectContentType(content)
	switch contentType {
	case
		"application/octet-stream",
		"application/x-tar",
		"application/x-bzip2",
		"application/x-gzip",
		"image/jpeg", "image/x-portable-pixmap",
		"image/x-ms-bmp", "image/x-icon", "image/svg+xml",
		"image/png", "image/gif", "image/x-xpmi",
		"application/postscript",
		"application/x-xz",
		"application/pdf",
		"application/pgp-signature",
		"application/zip":
		return false
	}
	return true
}

// IsFile returns true if the path is a regular file.
func IsFile(path string) bool {
	if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
		return true
	}
	return false
}

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "usage: %s <connection-uri> <db-name> <collection-name> <directory>\n", os.Args[0])
		os.Exit(1)
	}

	uri := os.Args[1]
	dbName := os.Args[2]
	colName := os.Args[3]
	searchDir := os.Args[4]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	col := db.Collection(colName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "path", Value: 1}, {Key: "line_count", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := col.Indexes().CreateOne(ctx, indexModel); err != nil {
		log.Fatalf("failed to create index: %v", err)
	}

	jobs := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go writeCode(ctx, col, jobs, &wg)

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsFile(path) {
			// Read first 512 bytes for content type detection
			data, err := os.ReadFile(path)
			if err != nil {
				return nil // skip unreadable files
			}
			if len(data) > 512 {
				data = data[:512]
			}
			if IsText(data) {
				jobs <- path
			}
		}
		return nil
	})
	close(jobs)
	wg.Wait()
}
