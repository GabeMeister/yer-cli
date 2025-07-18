import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

//
// HELPER FUNCTIONS
//

function handlePaintingBarCharts(evt) {
  const barChartElem = document.getElementById("container");
  const racingBarChartElem = document.getElementById(
    "racing-bar-chart-container"
  );
  const svgInPage = !!document.getElementsByTagName("svg").length;

  if (barChartElem !== null && !svgInPage) {
    paintBarChart();
  } else if (racingBarChartElem !== null && !svgInPage) {
    paintRacingBarChart();
  }
}

function handleChartJS(evt) {
  console.log("\n\n***** elem *****\n", evt.currentTarget, "\n\n");
  const dataz = evt.dataset.chartJsData;
  console.log("\n\n***** dataz *****\n", dataz, "\n\n");
  const ctx = document.getElementById("chart-js-canvas");
  if (ctx) {
    new Chart(ctx, {
      type: "bar",
      data: {
        labels: ["rb-frontend", "rb-backend", "demo-stack", "rb-docker"],
        datasets: [
          {
            label: "2024",
            data: [12, 19, 8, 15],
            backgroundColor: "rgba(54, 162, 235, 0.5)",
            borderColor: "rgba(54, 162, 235, 1)",
            borderWidth: 1,
          },
          {
            label: "2025",
            data: [8, 14, 12, 18],
            backgroundColor: "rgba(255, 99, 132, 0.5)",
            borderColor: "rgba(255, 99, 132, 1)",
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: "Lines of Code",
            },
          },
        },
        barPercentage: 1.0,
        categoryPercentage: 0.2,
      },
    });
  }
}

//
// HTMX LOAD
//

document.body.addEventListener("htmx:load", (elem) => {
  handlePaintingBarCharts(elem);
  handleChartJS(elem);
});

//
// KEYBOARD SHORTCUTS
//

window.addEventListener("keydown", (evt) => {
  if (evt.key === "ArrowLeft") {
    history.back();
  }
});
