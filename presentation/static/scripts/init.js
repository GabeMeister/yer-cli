import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

window.htmx.onLoad((elem) => {
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
});

window.addEventListener("keydown", (evt) => {
  if (evt.key === "ArrowLeft") {
    history.back();
  }
});
