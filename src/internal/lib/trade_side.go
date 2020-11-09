package lib

type TradeSide string

const (
	Buy  TradeSide = "buy"
	Sell TradeSide = "sell"
)

func (ts *TradeSide) IsBuy() bool {
	if *ts == Buy {
		return true
	}
	return false
}
