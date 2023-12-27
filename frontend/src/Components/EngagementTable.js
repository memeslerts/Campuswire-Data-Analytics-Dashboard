import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function EngagementTable({ day }) {
  const history = useNavigate();
  const [table, loadTable] = useState([]);
  const [dropdown, showDropdown] = useState(true);
  const [criteria, setCriteria] = useState("posts");

  useEffect(() => {
    axios
      .get(
        "http://localhost:8080/day" +
          day +
          "/dashboard/engagement?criteria=" +
          criteria +
          "&limit=10",
        { headers: { "Access-Control-Allow-Origin": "*" } }
      )
      .then((e) => {
        let arr = [];
        for (let j in e.data) {
          arr.push({
            name: e.data[j].author.firstName + " " + e.data[j].author.lastName,
            postCount: e.data[j].postCount,
            id: e.data[j].author.id,
          });
        }
        loadTable(arr);
      });
  }, [day, criteria]);

  return (
    <div className="h-full">
      <div className="dropdown w-full">
        <div
          tabIndex="0"
          role="button"
          className="btn w-full rounded-none rounded-t-lg"
          onClick={() => showDropdown(true)}
        >
          Filter
        </div>
        {dropdown && (
          <ul className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-full max-h-52 overflow-x-scroll flex flex-row">
            {[
              "posts",
              "comments",
              "postTime",
              "commentTime",
              "endorsedComments",
            ].map((e) => (
              <li className="w-full">
                <button
                  onClick={() => {
                    setCriteria(e);
                    showDropdown(false);
                  }}
                >
                  {e}
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
      <div className="text-center">
        <div className="mt-5 font-semibold">Engagement Table</div>
      </div>
      <div className="divider"></div>
      <table className="table overflow-scroll max-h-5">
        <thead>
          <th>#</th>
          <td>Name</td>
          <td>Posts</td>
        </thead>
        <tbody>
          {table.map((e, i) => {
            return (
              <tr className="hover" onClick={() => history("/User/" + e["id"] + "," + day)}>
                <th>{i + 1}</th>
                <td>{e["name"]}</td>
                <td>{e["postCount"]}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default EngagementTable;
