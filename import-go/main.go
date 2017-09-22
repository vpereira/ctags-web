package main

import (
	"bufio"
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

func IsText(content []byte) bool {
	contentType := http.DetectContentType(content)
	if contentType == "application/octet-stream" {
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

func ReadExtractFile(path string) []byte {
	fs, _ := os.Open(path)
	defer fs.Close()
	n := 512
	buff := make([]byte, n)
	fs.Read(buff)
	return buff
}

func main() {
	searchDir := "/home/vpereira/node-v6.11.1-linux-x64"
	jobs := make(chan string)
	session, err := mgo.Dial("172.17.0.2/ctags")
	if err != nil {
		panic(err)
	}
	c := session.DB("ctags").C("code")
	go writeCode(c, jobs)
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if IsFile(path) && IsText(ReadExtractFile(path)) {
			jobs <- path
		}
		return nil
	})
	close(jobs)
}
