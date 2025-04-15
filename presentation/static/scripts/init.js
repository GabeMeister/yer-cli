import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

//
// HELPER FUNCTIONS
//

function handleAnalyzeManuallyPageDragDrop(elem) {
  function updateHiddenInputs() {
    const left = [
      ...document
        .querySelector("#left")
        .querySelectorAll(".draggable-text")
        .values()
        .map((item) => item.textContent),
    ];

    const right = [
      ...document
        .querySelector("#right")
        .querySelectorAll(".draggable-text")
        .values()
        .map((item) => item.textContent),
    ];

    const leftHiddenInput = document.querySelector(
      "input[name='all-authors']"
    );
    leftHiddenInput.value = left.join(",");

    const rightHiddenInput = document.querySelector(
      "input[name='duplicate-authors']"
    );
    rightHiddenInput.value = right.join(",");

    // Handle the other form for actually submitting the duplicate
    const duplicateEngineersHiddenInput = document.querySelector(
      "input[name='duplicate-authors']"
    );
    duplicateEngineersHiddenInput.value = right.join(",");
  }

  let sortables = document.querySelectorAll(".sortable");

  for (let i = 0; i < sortables.length; i++) {
    let sortable = sortables[i];
    let sortableInstance = new Sortable(sortable, {
      animation: 150,

      group: "shared",

      ghostClass: "text-red-400",

      // Make the `.htmx-indicator` unsortable
      filter: ".htmx-indicator, .ignore-input",

      onMove: function (evt) {
        return evt.related.className.indexOf("htmx-indicator") === -1;
      },

      // Disable sorting on the `end` event
      onEnd: function (evt) {
        this.option("disabled", true);
        updateHiddenInputs();
        htmx.trigger("#shared-form", "submit");
      },
    });

    // Re-enable sorting on the `htmx:afterSwap` event
    sortable.addEventListener("htmx:afterSwap", function () {
      sortableInstance.option("disabled", false);
    });
  }
}

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
  handleAnalyzeManuallyPageDragDrop(elem);
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
