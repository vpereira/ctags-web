package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// WebContext holds data passed to the show.html template.
type WebContext struct {
	CodeLines []CodeLine
	FileName  string
	LineCount int
}

// Env holds the MongoDB client and current database/collection.
type Env struct {
	Client     *mongo.Client
	Db         *mongo.Database
	Collection *mongo.Collection
}

// OpenDB connects to MongoDB and selects the given database and collection.
func (env *Env) OpenDB(ctx context.Context, uri string, db string, collection string) error {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return err
	}
	env.Client = client
	env.Db = client.Database(db)
	env.Collection = env.Db.Collection(collection)
	return nil
}

func (env *Env) BrowsingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	showPath := filepath.Join("static", "show.html")
	tmpl, err := template.ParseFiles(showPath)
	if err != nil {
		log.Printf("failed to parse template: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	FilePath := r.FormValue("file")
	LineCount := 1
	if r.FormValue("linecount") != "" {
		LineCount, _ = strconv.Atoi(r.FormValue("linecount"))
	}

	// Query the code collection for all lines of the given file
	col := env.Db.Collection("code")
	cursor, err := col.Find(ctx, bson.M{"path": FilePath}, options.Find().SetSort(bson.D{{Key: "line_count", Value: 1}}))
	if err != nil {
		log.Printf("query error: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer cursor.Close(ctx)

	var results []CodeLine
	for cursor.Next(ctx) {
		var doc CodeLine
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("decode error: %v", err)
			continue
		}
		results = append(results, doc)
	}

	if len(results) == 0 {
		http.NotFound(w, r)
		return
	}

	context := WebContext{
		CodeLines: results,
		FileName:  FilePath,
		LineCount: LineCount,
	}

	if err := tmpl.Execute(w, context); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func (env *Env) TokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	token := r.FormValue("token")
	results, err := env.FindName(ctx, token)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	resultsSize := len(results)
	jsonResults, _ := json.Marshal(results)
	fmt.Fprintf(w, "{\"results\": %s,\"count\": %d}", jsonResults, resultsSize)
}

func (env *Env) FindName(ctx context.Context, name string) ([]Ctag, error) {
	var results []Ctag
	cursor, err := env.Collection.Find(ctx, bson.M{"name": bson.Regex{Pattern: name, Options: "i"}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc Ctag
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		results = append(results, doc)
	}
	return results, nil
}
