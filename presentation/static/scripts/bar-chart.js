import * as d3 from "https://cdn.jsdelivr.net/npm/d3@7/+esm";

/*
 * REQUIREMENTS: Must render html div with:
 * - id="container"
 * - data-value="[{"week_number": 1, "line_count": 1503}, ...]
 */

// Declare the chart dimensions and margins.
const width = 1400;
const height = 700;
const marginTop = 80;
const marginRight = 30;
const marginBottom = 100;
const marginLeft = 100;

const weekToMonth = {
  1: "Jan",
  2: "Jan",
  3: "Jan",
  4: "Jan",
  5: "Jan/Feb",
  6: "Feb",
  7: "Feb",
  8: "Feb",
  9: "Feb",
  10: "Mar",
  11: "Mar",
  12: "Mar",
  13: "Mar",
  14: "Mar/Apr",
  15: "Apr",
  16: "Apr",
  17: "Apr",
  18: "Apr",
  19: "May",
  20: "May",
  21: "May",
  22: "May",
  23: "May/Jun",
  24: "Jun",
  25: "Jun",
  26: "Jun",
  27: "Jun",
  28: "Jul",
  29: "Jul",
  30: "Jul",
  31: "Jul",
  32: "Jul/Aug",
  33: "Aug",
  34: "Aug",
  35: "Aug",
  36: "Aug",
  37: "Sep",
  38: "Sep",
  39: "Sep",
  40: "Sep",
  41: "Sep/Oct",
  42: "Oct",
  43: "Oct",
  44: "Oct",
  45: "Oct",
  46: "Nov",
  47: "Nov",
  48: "Nov",
  49: "Nov",
  50: "Nov/Dec",
  51: "Dec",
  52: "Dec",
};

function paintBarChart() {
  const elem = document.querySelector("#container");
  let data = JSON.parse(elem.getAttribute("data-value"));

  // Convert data to be chart-friendly
  data = data.map((d) => ({
    week_number: weekToMonth[d.week_number],
    line_count: d.line_count,
  }));

  /*
   * X SCALE CALCULATION
   */
  const x = d3
    .scaleBand()
    .domain(
      data.map(function (d) {
        return d.week_number;
      })
    )
    .range([marginLeft, width - marginRight])
    .padding(0.1);

  /*
   * Y SCALE CALCULATION
   */
  const y = d3
    .scaleLinear()
    .domain([0, d3.max(data, (d) => d.line_count)])
    .range([height - marginBottom, marginTop]);

  /*
   * INITIALIZE
   */
  const svg = d3
    .select("#container")
    .append("svg")
    .attr("width", width)
    .attr("height", height);

  /*
   * BACKGROUND
   */
  svg
    .append("g")
    .attr("id", "background")
    .append("rect")
    .attr("fill", "#f7f7f7")
    .attr("rx", "5px")
    .attr("ry", "5px")
    .attr("x", 0)
    .attr("y", 0)
    .attr("height", "100%")
    .attr("width", "100%");

  /*
   * X AXIS
   */
  svg
    .append("g")
    .attr("transform", `translate(0,${height - marginBottom})`)
    .attr("color", "black")
    .call(d3.axisBottom(x).tickSizeOuter(0))
    .selectAll("text")
    .attr("font-size", "18")
    .attr("dx", "2em")
    .attr("dy", "2rem")
    .attr("transform", "rotate(30)");

  /*
   * Y AXIS
   */
  const yAxis = svg
    .append("g")
    .attr("transform", `translate(${marginLeft},0)`)
    .attr("color", "black")
    .call(d3.axisLeft(y).tickFormat((y) => y.toLocaleString()))
    .call((g) => g.select(".domain").remove());
  yAxis.selectAll("text").attr("font-size", "18");
  yAxis.call((g) =>
    g
      .append("text")
      .attr("id", "y-axis-label")
      .attr("x", -marginLeft + 20)
      .attr("y", 40)
      .attr("font-size", "24")
      .attr("fill", "currentColor")
      .attr("text-anchor", "start")
      .text("↑ Line Count")
  );

  /*
   * BARS
   */
  svg
    .append("g")
    .attr("id", "bars")
    .attr("fill", "steelblue")
    .selectAll("rect")
    .data(data)
    .join("rect")
    .attr("x", (d) => x(d.week_number))
    .attr("y", (d) => {
      return y(0);
    })
    .attr("height", (_d) => {
      return 0;
    })
    .attr("width", x.bandwidth() - 10)
    .attr("rx", "5px")
    .attr("ry", "5px");

  /*
   * BAR ANIMATION
   */
  svg
    .selectAll("#bars rect")
    .transition()
    .ease(d3.easeBackOut)
    .duration(800)
    .attr("y", function (d) {
      return y(d.line_count);
    })
    .attr("height", function (d) {
      const zero = y(0);
      const yValue = y(d.line_count);

      return zero - yValue;
    })
    .delay(function (_d, i) {
      return i * 15;
    });
}

paintBarChart();
