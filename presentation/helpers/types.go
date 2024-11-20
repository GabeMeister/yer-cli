package presentation_helpers

type DataPoint struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

type FloatDataPoint struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

type BarChartData struct {
	Data       []DataPoint `json:"data"`
	XAxisLabel string      `json:"x_axis_label"`
	YAxisLabel string      `json:"y_axis_label"`
}

type BarChartFloatData struct {
	Data       []FloatDataPoint `json:"data"`
	XAxisLabel string           `json:"x_axis_label"`
	YAxisLabel string           `json:"y_axis_label"`
}
