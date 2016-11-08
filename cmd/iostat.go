package main

import (
	"github.com/thisisaaronland/go-iostat"
	"log"
)

func main() {

	stats, err := iostat.NewIOStatResults()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(stats)
}
