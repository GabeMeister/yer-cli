import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

function handleBarCharts(detail) {
  const barChartElem = document.getElementById("container");
  const svgInPage = !!document.getElementsByTagName("svg").length;

  if (barChartElem !== null && !svgInPage) {
    paintBarChart();
  }
}

function handleRacingBarCharts(detail) {
  const racingBarChartElem = document.getElementById(
    "racing-bar-chart-container"
  );
  const svgInPage = !!document.getElementsByTagName("svg").length;

  if (racingBarChartElem !== null && !svgInPage) {
    paintRacingBarChart();
  }
}

function handleGroupedBarCharts(evt) {
  const elem = evt.detail.elt;
  const dataElem = elem.querySelector("[data-grouped-bar-chart-data]");

  if (dataElem) {
    const jsonStr = dataElem.getAttribute("data-grouped-bar-chart-data");
    const chartData = JSON.parse(jsonStr);
    const canvasElem = elem.querySelector("#grouped-bar-chart-canvas");

    new Chart(canvasElem, chartData);
  }
}

function handleLineCharts(evt) {
  const elem = evt.detail.elt;
  const dataElem = elem.querySelector("[data-line-chart-data]");

  if (dataElem) {
    const jsonStr = dataElem.getAttribute("data-line-chart-data");
    const chartData = JSON.parse(jsonStr);
    const canvasElem = elem.querySelector("#line-chart-canvas");

    new Chart(canvasElem, chartData);
  }
}

//
// HTMX LOAD
//

document.body.addEventListener("htmx:load", (evt) => {
  handleBarCharts(evt);
  handleRacingBarCharts(evt);
  handleGroupedBarCharts(evt);
  handleLineCharts(evt);
});

//
// KEYBOARD SHORTCUTS
//

window.addEventListener("keydown", (evt) => {
  if (evt.key === "ArrowLeft") {
    history.back();
  }
});
