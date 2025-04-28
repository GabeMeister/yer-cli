import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

//
// HELPER FUNCTIONS
//

function handlePaintingBarCharts(elem) {
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

//
// HTMX LOAD
//

document.body.addEventListener("htmx:load", (elem) => {
  handlePaintingBarCharts(elem);
});

//
// KEYBOARD SHORTCUTS
//

window.addEventListener("keydown", (evt) => {
  if (evt.key === "ArrowLeft") {
    history.back();
  }
});
