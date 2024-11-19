package presentation_helpers

type DataPoint struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

type BarChartData struct {
	Data       []DataPoint `json:"data"`
	XAxisLabel string      `json:"x_axis_label"`
	YAxisLabel string      `json:"y_axis_label"`
}
