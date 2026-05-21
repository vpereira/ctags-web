package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
)

// Ctag represents a single universal-ctags JSON output record.
type Ctag struct {
	Type     string `bson:"_type" json:"type"`
	Name     string `bson:"name" json:"name"`
	Path     string `bson:"path" json:"path"`
	Pattern  string `bson:"pattern" json:"pattern"`
	Language string `bson:"language" json:"language"`
	Line     int    `bson:"line" json:"line"`
	Kind     string `bson:"kind" json:"kind"`
}

// CodeLine represents a single line of source code stored in the database.
type CodeLine struct {
	FilePath  string `bson:"path"`
	Line      string `bson:"line"`
	LineCount int    `bson:"line_count"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: %s <connection-uri> <db-name> <collection-name>\n", os.Args[0])
		os.Exit(1)
	}

	uri := os.Args[1]
	dbName := os.Args[2]
	colName := os.Args[3]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Opening db connection")
	env := &Env{}
	if err := env.OpenDB(ctx, uri, dbName, colName); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer env.Client.Disconnect(ctx)

	fmt.Println("Setting handlers")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/token", env.TokenHandler)
	http.HandleFunc("/show", env.BrowsingHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
