package components

type DataRangeType string

type DateRange struct {
	Options []DateRangeOption
}

type DateRangeOption struct {
	Name       DataRangeType
	Link       string
	IsSelected bool
}

var DateRangeValues = []DataRangeType{"24h", "1w", "1m", "3m", "1y"}
