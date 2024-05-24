package components

type DateRange struct {
	Options []DateRangeOption
}

type DateRangeOption struct {
	Name       string
	Link       string
	IsSelected bool
}

var DateRangeValues = []string{"24h", "1w", "1m", "3m", "1y"}
