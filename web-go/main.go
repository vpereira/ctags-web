package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Ctag struct {
	Type     string `json:"_type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Pattern  string `json:"pattern"`
	Language string `json:"language"`
	Line     int    `json:"line"`
	Kind     string `json:"kind"`
}

type CodeLine struct {
	FilePath  string `json:"path"`
	Line      string `json:"line"`
	LineCount int    `json:"line_count"`
}

type Env struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
}

type WebContext struct {
	CodeLines []CodeLine
	FileName  string
	LineCount int
}

func (env *Env) OpenDB(serverName string, db string, collection string) (*mgo.Session, error) {
	Session, err := mgo.Dial(serverName)
	env.Session = Session
	if err == nil {
		env.Session.SetMode(mgo.Monotonic, true)
	}
	env.Collection = env.SetDB(db, collection)
	return env.Session, err
}

func (env *Env) SetDB(db string, collection string) *mgo.Collection {
	return env.Session.DB(db).C(collection)
}

func (env *Env) BrowsingHandler(w http.ResponseWriter, r *http.Request) {
	var results []CodeLine

	//w.Header().Add("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("show.html")

	if err != nil {
		log.Fatal(err)
	}

	// it should be configurable
	col := env.SetDB("ctags", "code")
	FilePath := r.FormValue("file")
	LineCount := 0

	if r.FormValue("linecount") != "" {
		LineCount, _ = strconv.Atoi(r.FormValue("linecount"))
	}

	col.Find(bson.M{"filepath": FilePath}).Select(bson.M{"_id": 0, "linecount": 1, "line": 1}).Sort("linecount").All(&results)

	if len(results) <= 0 {
		http.Error(w, http.StatusText(404), 404)
	}

	context := WebContext{
		CodeLines: results,
		FileName:  FilePath,
		LineCount: LineCount,
	}

	tmpl.Execute(w, context)
}

func (env *Env) TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.FormValue("token")
	results, err := env.FindName(token)
	resultsSize := len(results)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	jsonResults, _ := json.Marshal(results)
	fmt.Fprintf(w, "{\"results\": %s,\"count\": %d}", jsonResults, resultsSize)
}

// must return the results
func (env *Env) FindName(name string) ([]Ctag, error) {
	var results []Ctag
	err := env.Collection.Find(bson.M{"name": bson.M{"$regex": bson.RegEx{name, ""}}}).All(&results)
	return results, err
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("usage: ", os.Args[0], " <server> <db-name> <collection-name>")
		os.Exit(1)
	}

	serverName := os.Args[1]
	dbName := os.Args[2]
	colName := os.Args[3]

	env := Env{}
	fmt.Println("Opening db connection")
	_, err := env.OpenDB(serverName, dbName, colName)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	defer env.Session.Close()
	fmt.Println("Setting Handlers")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/token", env.TokenHandler)
	http.HandleFunc("/show", env.BrowsingHandler)
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}
