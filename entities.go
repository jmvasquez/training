package main

import (
	"math"
	"math/rand"
	"strings"
)

type Price struct {
	Name    string  `json:"provider,omitempty"`
	DollarC float64 `json:"dollarC,omitempty"`
	DollarV float64 `json:"dollarV,omitempty"`
	EuroC   float64 `json:"euroC,omitempty"`
	EuroV   float64 `json:"euroV,omitempty"`
}

type Central struct {
	Price
}

type Bank struct {
	Price
}

type Broker struct {
	Price
}

type PriceFinder struct {
	Brokers []Broker
	Banks   []Bank
	Central Central
	//BankSource string
	Prices chan Price
}

type PriceProvider interface {
	getPrice(currency string, prices chan Price) func() error
}

func (b *Broker) getPrice(currency string, prices chan Price) func() error {
	return func() error {
		b.calculatePrices(banks)
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
		prices <- price //Ver si devolver Price o con el channel alcanza
		return nil
	}
}

func (b *Bank) getPrice(currency string, prices chan Price) func() error {
	return func() error {
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
		prices <- price //Ver si devolver Price o con el channel alcanza
		return nil
	}
}

func (b *Broker) calculatePrices(others []Bank) {
	var minDc, maxDc, minDv, maxDv, minEc, maxEc, minEv, maxEv float64

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
	var filename = "resources/banks.csv"

	//Loads banks data from CSV file
	loadBanksData(filename, &banks)

	//Initialization of statics brokers
	var balanz, bullExchange Broker
	balanz.Price.Name = "Balanz"
	bullExchange.Price.Name = "Bull Exchange"

	brokers = append(brokers, balanz, bullExchange)

	//Initialization of Central Bank
	central.Price = Price{
		Name:    "Banco Central",
		DollarC: 45.90,
		DollarV: 46.10,
		EuroC:   49.90,
		EuroV:   51.30,
	}

	//PriceFinder initialization   SOLE
	pfinder.Banks = banks
	pfinder.Brokers = brokers
	pfinder.Central = central
}
