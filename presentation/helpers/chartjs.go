package helpers

import (
	"GabeMeister/yer-cli/analyzer"
	"sort"
)

type YearComparisonChartData struct {
	Dataset    map[string]analyzer.YearComparison
	YAxisLabel string
}

func GetYearComparisonChartData(data YearComparisonChartData) map[string]interface{} {
	repos := []string{}
	prevYearData := []int{}
	currYearData := []int{}

	for repo := range data.Dataset {
		repos = append(repos, repo)
	}

	sort.Slice(repos, func(i, j int) bool {
		return data.Dataset[repos[i]].Curr > data.Dataset[repos[j]].Curr
	})

	for _, repo := range repos {
		item := data.Dataset[repo]
		prevYearData = append(prevYearData, item.Prev)
		currYearData = append(currYearData, item.Curr)
	}

	return map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": repos,
			"datasets": []map[string]interface{}{
				{
					"label":           "2024",
					"data":            prevYearData,
					"backgroundColor": "rgba(54, 162, 235, 0.5)",
					"borderColor":     "rgba(54, 162, 235, 1)",
					"borderWidth":     1,
				},
				{
					"label":           "2025",
					"data":            currYearData,
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
					},
				},
			},
			"barPercentage":      1.0,
			"categoryPercentage": 0.3,
		},
	}
}

// options: {
//   responsive: true,
//   scales: {
//     y: {
//       beginAtZero: true,
//       title: {
//         display: true,
//         text: "Lines of Code",
//       },
//     },
//   },
//   barPercentage: 1.0,
//   categoryPercentage: 0.2,
// },
// }
