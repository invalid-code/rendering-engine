package main

import (
	"log"
	"io"
	"os"
)

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	fileContent := string(data)
	return fileContent
}
