package candles

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"entrance/lection3/domain"
)

func FromPrices(prices <-chan domain.Price, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	go func() {
		cands := make(map[string]domain.Candle)
		for price := range prices {
			periodTS, err := domain.PeriodTS(period, price.TS)
			if err != nil {
				fmt.Println(err)
				continue
			}

			cand, ok := cands[price.Ticker]
			if !ok {
				cands[price.Ticker] = newCandleFromPrice(price, periodTS, period)
				continue
			}

			if cand.TS.Before(periodTS) {
				out <- cand
				cands[price.Ticker] = newCandleFromPrice(price, periodTS, period)
				continue
			}

			if price.Value > cand.High {
				cand.High = price.Value
			}
			if price.Value < cand.Low {
				cand.Low = price.Value
			}
			cand.Close = price.Value

			cands[price.Ticker] = cand
		}

		for _, cand := range cands {
			out <- cand
		}

		close(out)
	}()

	return out
}

func FromCandles(candles <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	go func() {
		cands := make(map[string]domain.Candle)
		for candIn := range candles {
			periodTS, err := domain.PeriodTS(period, candIn.TS)
			if err != nil {
				fmt.Println(err)
				continue
			}

			cand, ok := cands[candIn.Ticker]
			if !ok {
				candIn.Period = period
				candIn.TS = periodTS
				cands[candIn.Ticker] = candIn
				continue
			}

			if cand.TS.Before(periodTS) {
				out <- cand

				candIn.Period = period
				candIn.TS = periodTS
				cands[candIn.Ticker] = candIn
				continue
			}

			if candIn.High > cand.High {
				cand.High = candIn.High
			}
			if candIn.Low < cand.Low {
				cand.Low = candIn.Low
			}
			cand.Close = candIn.Close

			cands[candIn.Ticker] = cand
		}

		for _, cand := range cands {
			out <- cand
		}

		close(out)
	}()

	return out
}

func newCandleFromPrice(price domain.Price, periodTS time.Time, period domain.CandlePeriod) domain.Candle {
	return domain.Candle{
		Ticker: price.Ticker,
		Period: period,
		Open:   price.Value,
		High:   price.Value,
		Low:    price.Value,
		Close:  price.Value,
		TS:     periodTS,
	}
}

func ToCSV(fileName string, candles <-chan domain.Candle, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Ticker", "Period", "Open", "High", "Low", "Close", "Timestamp"}
	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for candle := range candles {
		// fmt.Println(candle)
		row := []string{
			candle.Ticker,
			string(candle.Period),
			strconv.FormatFloat(candle.Open, 'f', 2, 64),
			strconv.FormatFloat(candle.High, 'f', 2, 64),
			strconv.FormatFloat(candle.Low, 'f', 2, 64),
			strconv.FormatFloat(candle.Close, 'f', 2, 64),
			candle.TS.Format("2006-01-02T15:04:05-07:00"),
		}
		// fmt.Println(row)
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}
