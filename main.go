package main

import (
	"log"
	"net/http"
)

var banks []Bank
var brokers []Broker
var central Central
var pfinder PriceFinder

func main() {

	initStructs()

	//URL example: http://localhost:8000/?usd
	http.HandleFunc("/training", getBestPrice)

	//serve endpoints
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}

}
