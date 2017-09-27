package exchange

type Interface interface {
	Buy(float64) float64
	Sell(float64) float64
}

type Rate struct {
	Buy  float64
	Sell float64
}

func New(rate Rate) Interface {
	return &basic{Rate: rate}
}

type basic struct{ Rate Rate }

func (b *basic) Buy(x float64) float64  { return b.Rate.Buy * x }
func (b *basic) Sell(x float64) float64 { return b.Rate.Sell * x }
