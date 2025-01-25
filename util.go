package main

import (
	"io"
	"log"
	"math/rand"
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

func randFloats(min, max float32, n int) []float32 {
	res := make([]float32, n)
	for i := range res {
		res[i] = min + rand.Float32() * (max - min)
	}
	return res
}
