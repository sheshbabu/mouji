package components

type DateRange struct {
	Options []DateRangeOption
}

type DateRangeOption struct {
	Name       string
	Link       string
	IsSelected bool
}
