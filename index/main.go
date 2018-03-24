package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
{"_type": "tag", "name": " -valid (S)",
"path": "SUSE:SLE-10-SP3:Update:Test/samba/samba-3.0.36/docs/htmldocs/manpages/smb.conf.5.html",
"pattern": "/^<\\/h3><\\/div><\\/div><\\/div><a class=\"indexterm\" name=\"id2562560\"><\\/a><a name=\"-VALID\"><\\/a><div/",
"language": "HTML", "line": 5707, "kind": "heading3"}
*/
type Ctag struct {
	Type     string `json:"_type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Pattern  string `json:"pattern"`
	Language string `json:"language"`
	Line     int    `json:"line"`
	Kind     string `json:"kind"`
}

func main() {

	jobs := make(chan string)

	if len(os.Args) < 3 {
		fmt.Println("usage: ", os.Args[0], " <server> <json-file>")
		os.Exit(1)
	}

	serverName := os.Args[1]
	jsonFile := os.Args[2]

	file, err := os.Open(jsonFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	session, err := mgo.Dial(serverName)

	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	// TODO
	// configuration db and collection
	c := session.DB("ctags").C("ctags")

	go insertTags(c, jobs)

	reader := bufio.NewReader(file)

	for {
		l, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		jobs <- l
	}
	close(jobs)
}
func insertTags(col *mgo.Collection, jobs chan string) {
	for j := range jobs {
		bdoc := Ctag{}
		err := bson.UnmarshalJSON([]byte(j), &bdoc)
		if err != nil {
			log.Fatal(err)
		}
		// TODO
		// do it configure. Either create or update
		//_, err = col.Upsert(&bdoc, &bdoc)
		col.Insert(&bdoc)
		if err != nil {
			log.Fatal(err)
		}
	}
}
