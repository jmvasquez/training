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
	prices := make(chan Price, len(f.Banks) /*+len(f.Brokers)*/)
	g, _ := errgroup.WithContext(context.Background())

	//asi estaba antes
	// for _, p := range append(f.Brokers, f.banks...) {
	// 	g.Go(p.getPrice(currency, prices))
	// }

	// for _, p := range f.Brokers {
	// 	g.Go(p.getPrice(currency, prices))
	// }
	for _, j := range f.Banks {
		j := j
		g.Go(j.getPrice(currency, prices))
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
	close(prices)

	if len(prices) > 0 {
		pc := <-prices
		pv := pc
		for p := range prices {
			if currency == "usd" {
				if p.DollarC < pc.DollarC || pc.DollarC == 0 {
					pc = p
				}
				if p.DollarV > pc.DollarV {
					pv = p
				}
			} else {
				if p.EuroC < pc.EuroC {
					pc = p
				}
				if p.EuroV > pc.EuroV {
					pv = p
				}
			}
		}

		fmt.Printf("COMPRA: %v", pc)
		fmt.Printf("VENTA: %v", pv)
		cp <- pc
		cp <- pv
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

	select {
	case res := <-cp:
		jres, _ := json.Marshal(res)
		w.Write(jres)
		for p := range cp {
			jres, _ = json.Marshal(p)
			w.Write(jres)
		}
	case <-time.After(5 * time.Second):
		jcentral, _ := json.Marshal(pfinder.Central.Price)
		w.Write(jcentral)
	}
}
