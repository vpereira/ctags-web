package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
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


type WebContext struct {
	CodeLines []CodeLine
	FileName  string
	LineCount int
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
