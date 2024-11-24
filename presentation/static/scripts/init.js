import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

window.htmx.onLoad((elem) => {
  const barChartElem = document.getElementById("container");
  const racingBarChartElem = document.getElementById("bar-chart-container");

  if (barChartElem !== null) {
    paintBarChart();
  } else if (racingBarChartElem !== null) {
    paintRacingBarChart();
  }
});

window.addEventListener("keydown", (evt) => {
  if (evt.key === "ArrowLeft") {
    history.back();
  }
});
