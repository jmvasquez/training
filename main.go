package main

import (
	"log"
	"net/http"
)

var banks []Bank
var pfinder PriceFinder

func main() {

	initStructs()

	//URL example: http://localhost:8000/training?usd
	http.HandleFunc("/training", getBestPrice)

	//serve endpoints
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}

}
