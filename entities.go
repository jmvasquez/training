package main

import (
	"math"
	"math/rand"
	"strings"
)

//Price struct with currencies values
type Price struct {
	Name    string  `json:"provider,omitempty"`
	DollarC float64 `json:"dollar compra,omitempty"`
	DollarV float64 `json:"dollar venta,omitempty"`
	EuroC   float64 `json:"euro compra,omitempty"`
	EuroV   float64 `json:"euro venta,omitempty"`
}

//Bank struct
type Bank struct {
	Price
}

//Broker struct
type Broker struct {
	Price
}

//PriceFinder struct
type PriceFinder struct {
	Providers []PriceProvider
	Central   Bank
	Prices    chan Price
}

//PriceProvider interface that will calculate base on banks and brokers
type PriceProvider interface {
	getPrice(currency string, prices chan Price)
}

func getAllPrices(p PriceProvider, currency string, prices chan Price) func() error {
	return func() error {
		p.getPrice(currency, prices)
		return nil
	}
}

func (b Broker) getPrice(currency string, prices chan Price) {
	b.calculatePrices(banks, pfinder.Central)
	price := Price{}

	if strings.ToLower(currency) == "usd" {
		price.DollarC = b.Price.DollarC
		price.DollarV = b.Price.DollarV
		price.Name = b.Price.Name
	}
	if strings.ToLower(currency) == "euro" {
		price.EuroC = b.Price.EuroC
		price.EuroV = b.Price.EuroV
		price.Name = b.Price.Name
	}
	prices <- price
}

func (b Bank) getPrice(currency string, prices chan Price) {
	price := Price{}
	if strings.ToLower(currency) == "usd" {
		price.DollarC = b.Price.DollarC
		price.DollarV = b.Price.DollarV
		price.Name = b.Price.Name
	}
	if strings.ToLower(currency) == "euro" {
		price.EuroC = b.Price.EuroC
		price.EuroV = b.Price.EuroV
		price.Name = b.Price.Name
	}
	prices <- price
}

func (b *Broker) calculatePrices(others []Bank, central Bank) {
	var minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv float64

	minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv = getValues(central, minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv)
	for _, bank := range others {
		minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv = getValues(bank, minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv)
	}

	(*b).Price.DollarC = math.Round((minDc+rand.Float64()*(maxDc-minDc))*100) / 100
	(*b).Price.DollarV = math.Round((minDv+rand.Float64()*(maxDv-minDv))*100) / 100
	(*b).Price.EuroC = math.Round((minEc+rand.Float64()*(maxEc-minEc))*100) / 100
	(*b).Price.EuroV = math.Round((minEv+rand.Float64()*(maxEv-minEv))*100) / 100
}

func getValues(bank Bank, minDc float64, maxDc float64, minDv float64, maxDv float64, minEc float64, maxEc float64, minEv float64, maxEv float64) (float64, float64, float64, float64, float64, float64, float64, float64) {

	if minDc > bank.Price.DollarC || minDc < 0.1 {
		minDc = bank.Price.DollarC
	}
	if maxDc < bank.Price.DollarC || maxDc < 0.1 {
		maxDc = bank.Price.DollarC
	}
	if minDv > bank.Price.DollarV || minDv < 0.1 {
		minDv = bank.Price.DollarV
	}
	if maxDv < bank.Price.DollarV || maxDv < 0.1 {
		maxDv = bank.Price.DollarV
	}
	if minEc > bank.Price.EuroC || minEc < 0.1 {
		minEc = bank.Price.EuroC
	}
	if maxEc < bank.Price.EuroC || maxEc < 0.1 {
		maxEc = bank.Price.EuroC
	}
	if minEv > bank.Price.EuroV || minEv < 0.1 {
		minEv = bank.Price.EuroV
	}
	if maxEv < bank.Price.EuroV || maxEv < 0.1 {
		maxEv = bank.Price.EuroV
	}

	return minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv
}

func initStructs() {

	var brokers []Broker

	var filename = "resources/banks.csv"

	//Loads banks data from CSV file
	loadBanksData(filename, &banks)

	//Initialization of statics brokers
	var balanz, bullExchange Broker
	balanz.Price.Name = "Balanz"
	bullExchange.Price.Name = "Bull Exchange"

	brokers = append(brokers, balanz, bullExchange)

	//Initialization of Central Bank
	pfinder.Central.Price = Price{
		Name:    "Banco Central",
		DollarC: 45.90,
		DollarV: 46.10,
		EuroC:   49.90,
		EuroV:   51.30,
	}

	for _, bank := range banks {
		pfinder.Providers = append(pfinder.Providers, bank)
	}
	for _, broker := range brokers {
		pfinder.Providers = append(pfinder.Providers, broker)
	}
}
