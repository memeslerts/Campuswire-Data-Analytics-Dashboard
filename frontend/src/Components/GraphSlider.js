import { Bar, Pie, Line } from "react-chartjs-2";
import "chart.js/auto";
import axios from "axios";
import { useEffect, useState } from "react";

function GraphSlider({ sliderType, day }) {
  const [postCount, loadPostCount] = useState({ label: [], data: [] });
  const [answeredCount, loadAnsweredCount] = useState([0, 0]);
  const [unreadPosts, loadUnreadPosts] = useState({ label: [], data: [] });
  const [viewCount, loadViewCount] = useState({ label: [], data: [] });

  useEffect(() => {
    axios
      .get(
        "http://localhost:8080/day" + day + "/dashboard/graph/postCount?range=7",
        { headers: { "Access-Control-Allow-Origin": "*" } }
      )
      .then((e) => {
        loadPostCount({ label: e.data.date, data: e.data.posts });
      });

    axios
      .get("http://localhost:8080/day" + day + "/dashboard/graph/answerData", {
        headers: { "Access-Control-Allow-Origin": "*" },
      })
      .then((e) => {
        loadAnsweredCount([e.data.AnsweredPosts, e.data.UnansweredPosts]);
      });

    axios
      .get(
        "http://localhost:8080/day" + day + "/dashboard/graph/unread?range=7",
        {
          headers: { "Access-Control-Allow-Origin": "*" },
        }
      )
      .then((e) => {
        var label = [];
        var data = [];
        for (const [key, value] of Object.entries(e.data)) {
          label.push(key.substring(5));
          data.push(value["UnreadPosts"].length);
        }
        loadUnreadPosts({ label: label, data: data });
      });
      axios
      .get(
        "http://localhost:8080/day" + day + "/dashboard/graph/viewCount",
        {
          headers: { "Access-Control-Allow-Origin": "*" },
        }
      )
      .then((e) => {
        loadViewCount({ label: e.data.date, data: e.data.views });
      });
  }, [day]);

  let pie_graph = [
    <Pie
      data={{
        labels: ["Answered", "Unanswered"],
        datasets: [
          {
            label: "Answered Posts",
            data: answeredCount,
            backgroundColor: ["rgb(255, 99, 132)", "rgb(54, 162, 235)"],
            hoverOffset: 4,
          },
        ],
      }}
    />,
  ];

  let bar_graph = [
    <Bar
      data={{
        labels: postCount.label,
        datasets: [
          {
            label: "Post Count",
            data: postCount.data,
            backgroundColor: [
              "rgba(255, 99, 132, 0.2)",
              "rgba(255, 159, 64, 0.2)",
              "rgba(255, 205, 86, 0.2)",
              "rgba(75, 192, 192, 0.2)",
              "rgba(54, 162, 235, 0.2)",
              "rgba(153, 102, 255, 0.2)",
              "rgba(201, 203, 207, 0.2)",
            ],
            borderColor: [
              "rgb(255, 99, 132)",
              "rgb(255, 159, 64)",
              "rgb(255, 205, 86)",
              "rgb(75, 192, 192)",
              "rgb(54, 162, 235)",
              "rgb(153, 102, 255)",
              "rgb(201, 203, 207)",
            ],
            borderWidth: 1,
          },
        ],
      }}
      width={"700px"}
      options={{ maintainAspectRatio: false }}
    />,
    <Line
      data={{
        labels: unreadPosts.label,
        datasets: [
          {
            label: "Unread Posts",
            data: unreadPosts.data,
            fill: false,
            borderColor: "rgb(75, 192, 192)",
            tension: 0.1,
          },
        ],
      }}
      width={"700px"}
      options={{ maintainAspectRatio: false }}
    />,<Line
    data={{
      labels: viewCount.label,
      datasets: [
        {
          label: "View Counts",
          data: viewCount.data,
          fill: false,
          borderColor: "rgb(239 68 68)",
          tension: 0.1,
        },
      ],
    }}
    width={"700px"}
    options={{ maintainAspectRatio: false }}
  />
  ];

  let graphs = sliderType == "Pie" ? pie_graph : bar_graph;
  return (
    <div className="outline outline-4 outline-gray-300 rounded-lg flex flex-col items-stretch shadow-xl">
      <div className="carousel grow" style={{maxWidth: '700px'}}>
        {graphs.map((e) => {
          return (
            <div className="carousel-item relative w-full">
              <span className="">{e}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default GraphSlider;
