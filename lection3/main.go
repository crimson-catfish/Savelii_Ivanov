package main

import (
	"context"
	"time"

	"entrance/lection3/candles"
	"entrance/lection3/domain"
	"entrance/lection3/generator"
)

func main() {
	cfg := generator.Config{
		Factor:  500,
		Delay:   time.Millisecond * 10,
		Tickers: []string{"NVID", "AAPL", "SBER"},
	}
	pg := generator.NewPricesGenerator(cfg)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	priceChan := pg.Prices(ctx)

	oneMinuteCandleChan := candles.FromPrices(priceChan, domain.CandlePeriod1m)
	twoMinuteCandleChan := candles.FromCandles(oneMinuteCandleChan, domain.CandlePeriod2m)
	tenMinuteCandleChan := candles.FromCandles(twoMinuteCandleChan, domain.CandlePeriod10m)

}
