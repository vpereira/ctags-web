package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	http.ListenAndServe(":8080", nil)
}
