package components

type BarChartInputDataPoint struct {
	Label string
	Data  int
}

type BarChart struct {
	Height float64
	Width  float64
	Data   []BarChartDataPoint
}

type BarChartDataPoint struct {
	X         float64
	Y         float64
	Width     float64
	Height    float64
	MaxHeight float64
	TopOffset float64
	Value     int
	Label     string
}

// SVG Coordinate System: https://developer.mozilla.org/en-US/docs/Web/SVG/Tutorial/Positions#the_grid
func NewBarChart(input []BarChartInputDataPoint) BarChart {
	chartWidth := 900.0
	chartHeight := 200.0
	chartTopOffset := 20.0
	chartLeftOffset := 0.0

	availableBarHeight := chartHeight - chartTopOffset
	availableBarWidth := (chartWidth - chartLeftOffset) / float64(len(input))
	barWidth := availableBarWidth / 2
	barGutterWidth := barWidth

	maxValue := 0
	for _, dataPoint := range input {
		if dataPoint.Data > maxValue {
			maxValue = dataPoint.Data
		}
	}

	barHeightScaleFactor := float64(availableBarHeight) / float64(maxValue)

	bars := []BarChartDataPoint{}
	for i, dataPoint := range input {
		bar := BarChartDataPoint{
			X:         chartLeftOffset + (float64(i) * availableBarWidth) + (barGutterWidth / 2),
			Width:     barWidth,
			Y:         chartTopOffset + availableBarHeight - (float64(dataPoint.Data) * barHeightScaleFactor),
			Height:    chartTopOffset + float64(dataPoint.Data)*barHeightScaleFactor,
			MaxHeight: chartTopOffset + availableBarHeight,
			TopOffset: chartTopOffset,
			Value:     dataPoint.Data,
			Label:     dataPoint.Label,
		}
		bars = append(bars, bar)
	}

	return BarChart{
		Height: chartHeight,
		Width:  chartWidth,
		Data:   bars,
	}
}
