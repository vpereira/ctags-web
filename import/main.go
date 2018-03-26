package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	mgo "gopkg.in/mgo.v2"
)

type CodeLine struct {
	FilePath  string `json:"path"`
	Line      string `json:"line"`
	LineCount int    `json:"line_count"`
}

func readFile(fileName string) []CodeLine {
	var codeLines []CodeLine
	lineCount := 1
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		doc := CodeLine{LineCount: lineCount, Line: scanner.Text(), FilePath: fileName}
		codeLines = append(codeLines, doc)
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return codeLines
}

func writeCode(c *mgo.Collection, jobs chan string) {
	for j := range jobs {
		for _, doc := range readFile(j) {
			err := c.Insert(&doc)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// maybe just match if content-type ~= /^text/
// and application/xml and inode/x-empty
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

func IsFile(path string) bool {
	if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
		return true
	}
	return false
}

// extract the first 512 bytes from file
// use it for the detection function IsText
// the 512 should be configurable
func ReadExtractFile(path string) []byte {
	fs, _ := os.Open(path)
	defer fs.Close()
	n := 512
	buff := make([]byte, n)
	fs.Read(buff)
	return buff
}

func main() {

	if len(os.Args) < 5 {
		fmt.Println("usage: ", os.Args[0], " <server> <db-name> <collection-name> <directory>")
		os.Exit(1)
	}

	serverName := os.Args[1]
	dbName := os.Args[2]
	colName := os.Args[3]
	searchDir := os.Args[4]

	jobs := make(chan string)
	session, err := mgo.Dial(serverName)
	if err != nil {
		panic(err)
	}
	c := session.DB(dbName).C(colName)

	go writeCode(c, jobs)
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if IsFile(path) && IsText(ReadExtractFile(path)) {
			jobs <- path
		}
		return nil
	})
	close(jobs)
}
