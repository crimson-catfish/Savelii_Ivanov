package main

import (
	"context"
	"sync"
	"time"

	"entrance/lection3/candles"
	"entrance/lection3/domain"
	"entrance/lection3/generator"
)

func main() {
	cfg := generator.Config{
		Factor:  80,
		Delay:   time.Nanosecond,
		Tickers: []string{"NVID", "AAPL", "SBER"},
	}
	pg := generator.NewPricesGenerator(cfg)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	priceChan := pg.Prices(ctx)

	wg := sync.WaitGroup{}

	oneMinuteCandleChan := candles.FromPrices(priceChan, domain.CandlePeriod1m)
	oneMinuteCandleChan1 := make(chan domain.Candle)
	oneMinuteCandleChan2 := make(chan domain.Candle)
	broadcast(oneMinuteCandleChan, oneMinuteCandleChan1, oneMinuteCandleChan2)
	wg.Add(1)
	go candles.ToCSV("lection3/candles_1m.csv", oneMinuteCandleChan1, &wg)

	twoMinuteCandleChan := candles.FromCandles(oneMinuteCandleChan2, domain.CandlePeriod2m)
	twoMinuteCandleChan1 := make(chan domain.Candle)
	twoMinuteCandleChan2 := make(chan domain.Candle)
	broadcast(twoMinuteCandleChan, twoMinuteCandleChan1, twoMinuteCandleChan2)
	wg.Add(1)
	go candles.ToCSV("lection3/candles_2m.csv", twoMinuteCandleChan1, &wg)

	tenMinuteCandleChan := candles.FromCandles(twoMinuteCandleChan2, domain.CandlePeriod10m)
	wg.Add(1)
	go candles.ToCSV("lection3/candles_10m.csv", tenMinuteCandleChan, &wg)

	wg.Wait()
}

func broadcast[T any](input <-chan T, outputs ...chan T) {
	go func() {
		for val := range input {
			for _, out := range outputs {
				out <- val
			}
		}
		for _, out := range outputs {
			close(out)
		}
	}()
}
