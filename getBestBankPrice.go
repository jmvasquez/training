package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func cleanPrices(pc *Price, pv *Price, currency string) {
	(*pc).Name = "Mejor Precio Compra: " + (*pc).Name
	(*pv).Name = "Mejor Precio Venta: " + (*pv).Name
	if currency == "usd" {
		(*pc).DollarV = 0.0
		(*pc).EuroC = 0.0
		(*pc).EuroV = 0.0
		(*pv).DollarC = 0.0
		(*pv).EuroC = 0.0
		(*pv).EuroV = 0.0
	} else {
		(*pc).DollarC = 0.0
		(*pc).DollarV = 0.0
		(*pc).EuroV = 0.0
		(*pv).DollarC = 0.0
		(*pv).DollarV = 0.0
		(*pv).EuroC = 0.0

	}
}

func (f *PriceFinder) findBestPrices(currency string, cp chan Price) {
	prices := make(chan Price, len(f.Providers))
	g, _ := errgroup.WithContext(context.Background())

	for _, p := range f.Providers {
		p := p
		g.Go(getAllPrices(p, currency, prices))
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
	close(prices)

	if len(prices) > 0 {
		pc := Price{}
		pv := Price{}
		for p := range prices {
			if currency == "usd" {
				if p.DollarC < pc.DollarC || pc.DollarC == 0.0 {
					pc = p
				}
				if p.DollarV > pv.DollarV || pv.DollarC == 0.0 {
					pv = p
				}
			} else {
				if p.EuroC < pc.EuroC || pc.EuroC == 0.0 {
					pc = p
				}
				if p.EuroV > pv.EuroV || pv.EuroC == 0.0 {
					pv = p
				}
			}
		}

		cleanPrices(&pc, &pv, currency)
		cp <- pv
		cp <- pc
	}
	close(cp)
}

func getBestPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("Error, GET is the only supported method \n")
		http.Error(w, "unsupported method", http.StatusBadRequest)
		return
	}

	params := strings.Split(r.URL.RawQuery, "&")
	currency := strings.ToLower(params[0])

	if currency != "usd" && currency != "euro" {
		log.Printf("Error, invalid parameter. \n")
		http.Error(w, "Invalid parameter. Should be \"usd\" or \"euro\"", http.StatusBadRequest)
		return
	}

	cp := make(chan Price, 2)
	go pfinder.findBestPrices(currency, cp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var results []Price
	select {
	case res := <-cp:
		results = append(results, res)
		for p := range cp {
			results = append(results, p)
		}
		jres, _ := json.Marshal(results)
		w.Write(jres)
	case <-time.After(5 * time.Millisecond):
		jcentral, _ := json.Marshal(pfinder.Central)
		w.Write(jcentral)
	}
}
