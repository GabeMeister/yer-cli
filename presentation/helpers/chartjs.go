package helpers

func GetGroupedBarChartData(datasetNames []string, values [][]int) map[string]interface{} {
	return map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": []string{"rb-frontend", "rb-backend", "demo-stack", "rb-docker"},
			"datasets": []map[string]interface{}{
				{
					"label":           "2024",
					"data":            []int{12, 19, 8, 15},
					"backgroundColor": "rgba(54, 162, 235, 0.5)",
					"borderColor":     "rgba(54, 162, 235, 1)",
					"borderWidth":     1,
				},
				{
					"label":           "2025",
					"data":            []int{8, 14, 12, 18},
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
						"text":    "Lines of Code",
					},
				},
			},
			"barPercentage":      1.0,
			"categoryPercentage": 0.2,
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
