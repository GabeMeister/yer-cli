import { paintBarChart } from "./bar-chart.js";
import { paintRacingBarChart } from "./racing-bar-chart.js";

function handleSortables(elem) {
  function updateHiddenInputs() {
    const left = [
      ...document
        .querySelector("#left")
        .querySelectorAll("input[name='item']")
        .values()
        .map((item) => item.value),
    ];

    const right = [
      ...document
        .querySelector("#right")
        .querySelectorAll("input[name='item']")
        .values()
        .map((item) => item.value),
    ];

    const leftHiddenInput = document.querySelector(
      "input[name='left-form-items']"
    );
    leftHiddenInput.value = left.join(",");

    const rightHiddenInput = document.querySelector(
      "input[name='right-form-items']"
    );
    rightHiddenInput.value = right.join(",");
  }

  let sortables = elem.querySelectorAll(".sortable");

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

window.htmx.onLoad((elem) => {
  handleSortables(elem);

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
