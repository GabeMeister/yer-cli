package helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"fmt"
	"sort"
)

type YearComparisonChartData struct {
	Dataset    map[string]analyzer.YearComparison
	YAxisLabel string
}

type BarChartItem struct {
	Name  string
	Value int
}

type ChartJSBarChartData struct {
	Dataset    []BarChartItem
	YAxisLabel string
	XAxisLabel string
}

type LineChartDataset struct {
	Name    string
	Dataset []int
}

type LineChartData struct {
	Datasets        []LineChartDataset
	XAxisTickLabels []string
	YAxisLabel      string
	XAxisLabel      string
}

type JSObject map[string]interface{}

func GetMonthsOfYear() []string {
	return []string{
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Dec",
	}
}

func GetWeekDays() []string {
	return []string{
		"Sun",
		"Mon",
		"Tues",
		"Wed",
		"Thur",
		"Fri",
		"Sat",
	}
}

func GetHoursOfDay() []string {
	return []string{
		"12am",
		"1am",
		"2am",
		"3am",
		"4am",
		"5am",
		"6am",
		"7am",
		"8am",
		"9am",
		"10am",
		"11am",
		"12pm",
		"1pm",
		"2pm",
		"3pm",
		"4pm",
		"5pm",
		"6pm",
		"7pm",
		"8pm",
		"9pm",
		"10pm",
		"11pm",
	}
}

// For each week of the year, get the month abbreviation that the week falls in
func GetMonthsThroughYear() []string {
	nums := []string{}
	for i := 1; i <= 52; i++ {
		nums = append(nums, WeekToMonth(i))
	}

	return nums
}

func WeekToMonth(weekNum int) string {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	// Approximate week to month (4.33 weeks per month on average)
	monthIndex := (weekNum - 1) / 4
	if monthIndex > 11 {
		monthIndex = 11
	}
	if monthIndex < 0 {
		monthIndex = 0
	}

	return months[monthIndex]
}

type BarChartOptions struct {
	Sort bool
}

func GetBarChartData(data ChartJSBarChartData, options BarChartOptions) map[string]interface{} {
	dataset := data.Dataset

	if options.Sort {
		sort.Slice(dataset, func(i int, j int) bool {
			return dataset[i].Value > dataset[j].Value
		})
	}

	labels := []string{}
	values := []int{}
	for _, item := range dataset {
		labels = append(labels, item.Name)
		values = append(values, item.Value)
	}

	return map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": labels,
			"datasets": []map[string]interface{}{
				{
					"label":           "2025",
					"data":            values,
					"backgroundColor": "rgba(255, 99, 132, 0.5)",
					"borderColor":     "rgba(255, 99, 132, 1)",
					"borderWidth":     1,
				},
			},
		},
		"options": map[string]interface{}{
			"responsive": true,
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
					"title": map[string]interface{}{
						"display": true,
						"text":    data.YAxisLabel,
						"color":   "white",
					},
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
				"x": map[string]interface{}{
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
			},
			"plugins": map[string]interface{}{
				"legend": map[string]interface{}{
					"labels": map[string]interface{}{
						"color": "white",
					},
				},
			},
			"barPercentage":      1.0,
			"categoryPercentage": 0.5,
		},
	}
}

type YearComparisonOptions struct {
	Limit int
}

func GetYearComparisonChartData(data YearComparisonChartData, options YearComparisonOptions) map[string]interface{} {
	// Could be repos, authors, etc.
	buckets := []string{}

	prevYearData := []int{}
	currYearData := []int{}

	for repo := range data.Dataset {
		buckets = append(buckets, repo)
	}

	sort.Slice(buckets, func(i, j int) bool {
		return data.Dataset[buckets[i]].Curr > data.Dataset[buckets[j]].Curr
	})

	if options.Limit != 0 && len(data.Dataset) > options.Limit {
		buckets = buckets[0:options.Limit]
	}

	for _, repo := range buckets {
		item := data.Dataset[repo]
		prevYearData = append(prevYearData, item.Prev)
		currYearData = append(currYearData, item.Curr)
	}

	return map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": buckets,
			"datasets": []map[string]interface{}{
				{
					"label":           "2024",
					"data":            prevYearData,
					"backgroundColor": "rgba(54, 162, 235, 0.5)",
					"borderColor":     "rgba(54, 162, 235, 1)",
					"borderWidth":     1,
					"stack":           "Stack 1",
				},
				{
					"label":           "2025",
					"data":            currYearData,
					"backgroundColor": "rgba(255, 99, 132, 0.5)",
					"borderColor":     "rgba(255, 99, 132, 1)",
					"borderWidth":     1,
					"stack":           "Stack 0",
				},
			},
		},
		"options": map[string]interface{}{
			"responsive": true,
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
					"title": map[string]interface{}{
						"display": true,
						"text":    data.YAxisLabel,
						"color":   "white",
					},
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
				"x": map[string]interface{}{
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
			},
			"plugins": map[string]interface{}{
				"legend": map[string]interface{}{
					"labels": map[string]interface{}{
						"color": "white",
					},
				},
			},
			"barPercentage":      1.0,
			"categoryPercentage": 0.3,
		},
	}
}

type StackedBarChartDataset struct {
	// The thing that the key shows
	Label string
	// The actual numbers across all the items on the X axis
	Data []int
	// The background color
	BackgroundColor string
	// The "stack" that this dataset belongs to within the bucket. (Used to
	// "stack" related things together within one line)
	Stack string
}

// func GetStackedBarChartData(datasets []StackedBarChartDataset) map[string]interface{} {
// 	// Could be repos, authors, etc.
// 	buckets := []string{}

// 	prevYearData := []int{}
// 	currYearData := []int{}

// 	for repo := range data.Dataset {
// 		buckets = append(buckets, repo)
// 	}

// 	sort.Slice(buckets, func(i, j int) bool {
// 		return data.Dataset[buckets[i]].Curr > data.Dataset[buckets[j]].Curr
// 	})

// 	if options.Limit != 0 && len(data.Dataset) > options.Limit {
// 		buckets = buckets[0:options.Limit]
// 	}

// 	for _, repo := range buckets {
// 		item := data.Dataset[repo]
// 		prevYearData = append(prevYearData, item.Prev)
// 		currYearData = append(currYearData, item.Curr)
// 	}

// 	return map[string]interface{}{
// 		"type": "bar",
// 		"data": map[string]interface{}{
// 			"labels": buckets,
// 			"datasets": []map[string]interface{}{
// 				{
// 					"label":           "2024",
// 					"data":            prevYearData,
// 					"backgroundColor": "rgba(54, 162, 235, 0.5)",
// 					"borderColor":     "rgba(54, 162, 235, 1)",
// 					"borderWidth":     1,
// 					"stack":           "Stack 1",
// 				},
// 				{
// 					"label":           "2025",
// 					"data":            currYearData,
// 					"backgroundColor": "rgba(255, 99, 132, 0.5)",
// 					"borderColor":     "rgba(255, 99, 132, 1)",
// 					"borderWidth":     1,
// 					"stack":           "Stack 0",
// 				},
// 			},
// 		},
// 		"options": map[string]interface{}{
// 			"responsive": true,
// 			"scales": map[string]interface{}{
// 				"y": map[string]interface{}{
// 					"beginAtZero": true,
// 					"title": map[string]interface{}{
// 						"display": true,
// 						"text":    data.YAxisLabel,
// 						"color":   "white",
// 					},
// 					"ticks": map[string]interface{}{
// 						"color": "white",
// 					},
// 				},
// 				"x": map[string]interface{}{
// 					"ticks": map[string]interface{}{
// 						"color": "white",
// 					},
// 				},
// 			},
// 			"plugins": map[string]interface{}{
// 				"legend": map[string]interface{}{
// 					"labels": map[string]interface{}{
// 						"color": "white",
// 					},
// 				},
// 			},
// 			"barPercentage":      1.0,
// 			"categoryPercentage": 0.3,
// 		},
// 	}
// }

var COLORS = []string{
	"rgb(255, 99, 132)",  // Red
	"rgb(54, 162, 235)",  // Blue
	"rgb(255, 205, 86)",  // Yellow
	"rgb(75, 192, 192)",  // Teal
	"rgb(153, 102, 255)", // Purple
	"rgb(255, 159, 64)",  // Orange
	"rgb(199, 199, 199)", // Grey
	"rgb(83, 102, 255)",  // Indigo
	"rgb(255, 99, 255)",  // Magenta
	"rgb(99, 255, 132)",  // Green
	"rgb(255, 192, 203)", // Pink
	"rgb(173, 216, 230)", // Light Blue
	"rgb(144, 238, 144)", // Light Green
	"rgb(255, 218, 185)", // Peach
	"rgb(221, 160, 221)", // Plum
	"rgb(255, 228, 181)", // Moccasin
	"rgb(176, 196, 222)", // Light Steel Blue
	"rgb(255, 182, 193)", // Light Pink
	"rgb(152, 251, 152)", // Pale Green
	"rgb(255, 160, 122)", // Light Salmon
}

func GetLineChartData(data LineChartData) map[string]interface{} {
	if len(data.XAxisTickLabels) != len(data.Datasets[0].Dataset) {
		panic(fmt.Sprintf("Incorrect input for line chart data! %+v", data))
	}

	datasets := []JSObject{}
	for i, item := range data.Datasets {
		datasets = append(datasets, JSObject{
			"label":       item.Name,
			"data":        item.Dataset,
			"fill":        false,
			"borderColor": COLORS[i],
			"tension":     0.1,
		})
	}

	return map[string]interface{}{
		"type": "line",
		"data": map[string]interface{}{
			"labels":   data.XAxisTickLabels,
			"datasets": datasets,
		},
		"options": map[string]interface{}{
			"responsive": true,
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
					"title": map[string]interface{}{
						"display": true,
						"text":    data.YAxisLabel,
						"color":   "white",
					},
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
				"x": map[string]interface{}{
					"title": map[string]interface{}{
						"display": true,
						"text":    data.XAxisLabel,
						"color":   "white",
					},
					"ticks": map[string]interface{}{
						"color": "white",
					},
				},
			},
			"plugins": map[string]interface{}{
				"legend": map[string]interface{}{
					"labels": map[string]interface{}{
						"color": "white",
					},
				},
			},
		},
	}
}
