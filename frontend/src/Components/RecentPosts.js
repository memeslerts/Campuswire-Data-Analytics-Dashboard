import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function RecentPosts({ day }) {
  const history = useNavigate();
  const [table, loadTable] = useState([]);
  const [criteria, setCriteria] = useState("posts");

  useEffect(() => {
    axios
      .get(
        "http://localhost:8080/day" +
          day +
          "/posts?limit=5&reverse=true&criteria=publishedAt",
        { headers: { "Access-Control-Allow-Origin": "*" } }
      )
      .then((e) => {
        let arr = [];
        for (let j in e.data) {
          arr.push({
            name: e.data[j].author.firstName + " " + e.data[j].author.lastName,
            title: e.data[j].title,
            answersCount: e.data[j].answersCount,
            viewsCount: e.data[j].viewsCount,
            id: e.data[j].id
          });
        }
        loadTable(arr);
      });
  }, [day, criteria]);

  return (
    <div className="w-full">
      <div className="font-semibold text-center p-4">Recent Posts</div>
      <div className="divider m-0"></div>
      <table className="table h-1 overflow-scroll">
        <thead>
            <tr>
                <td>Name</td>
                <td>Title</td>
                <td>Answers</td>
                <td>Views</td>
            </tr>
        </thead>
        <tbody className="h-full overflow-y-scroll">
          {table.map((e) => {
            return (
              <tr className="hover" onClick={() => history("/Post/" + e["id"])}>
                <td>{e["name"]}</td>
                <td>{e["title"]}</td>
                <td>{e["answersCount"]}</td>
                <td>{e["viewsCount"]}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default RecentPosts;