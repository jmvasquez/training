package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

//Load data from csv file and store it into the struct []Bank
func loadBanksData(filename string, data *[]Bank) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Fail to open csv file. %s", err)
		os.Exit(1)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalf("Fail to read lines from csv file. %s", err)
		os.Exit(1)
	}

	var unit Bank
	for _, line := range lines {
		name := line[0]
		dollarC, err := strconv.ParseFloat(line[1], 64)
		dollarV, err := strconv.ParseFloat(line[2], 64)
		euroC, err := strconv.ParseFloat(line[3], 64)
		euroV, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			log.Fatalf("Fail to parse %s", err)
		}
		unit.Price = Price{
			Name:    name,
			DollarC: dollarC,
			DollarV: dollarV,
			EuroC:   euroC,
			EuroV:   euroV,
		}

		(*data) = append((*data), unit)
	}
}
