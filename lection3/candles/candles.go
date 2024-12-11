package candles

import (
	"fmt"
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
