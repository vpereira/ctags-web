package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
  "path/filepath"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Env struct {
	Session    *mgo.Session
	Db         *mgo.Database
	Collection *mgo.Collection
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

	wd, _ := os.Getwd()

	println("working dir", wd)

  showPath := filepath.Join(wd,"static", "show.html")

  println("show path", showPath)

	tmpl, err := template.ParseFiles(showPath)

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
