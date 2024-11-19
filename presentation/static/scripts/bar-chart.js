import * as d3 from "https://cdn.jsdelivr.net/npm/d3@7/+esm";

/*
 * REQUIREMENTS: Must render html div with id="container" and the following in it's data-value attribute something like this:

{
    "data": [
      {"x": "My Label 1", "y": 1503}, 
      {"x": "My Label 2", "y": 1625} 
    ],
    "x_axis_label": "Commits",
    "y_axis_label": "Line Count"
}
 
 */

const width = 1400;
const height = 700;
const marginTop = 80;
const marginRight = 30;
const marginBottom = 100;
const marginLeft = 100;

function paintBarChart() {
  const elem = document.querySelector("#container");
  let chartData = JSON.parse(elem.getAttribute("data-value"));
  let data = chartData.data;
  let yAxisLabel = chartData.y_axis_label;

  /*
   * X SCALE CALCULATION
   */
  const x = d3
    .scaleBand()
    .domain(
      data.map(function (d) {
        return d.x;
      })
    )
    .range([marginLeft, width - marginRight])
    .padding(0.1);

  /*
   * Y SCALE CALCULATION
   */
  const y = d3
    .scaleLinear()
    .domain([0, d3.max(data, (d) => d.y)])
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
      .text(yAxisLabel)
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
    .attr("x", (d) => x(d.x))
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
      return y(d.y);
    })
    .attr("height", function (d) {
      const zero = y(0);
      const yValue = y(d.y);

      return zero - yValue;
    })
    .delay(function (_d, i) {
      return i * 15;
    });
}

paintBarChart();
