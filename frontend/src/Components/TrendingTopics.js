import { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function TendingTopcis({ day }) {
  const [words, setWords] = useState([]);
  const history = useNavigate();

  useEffect(() => {
    axios
      .get("http://localhost:8080/day" + day + "/dashboard/topics", {
        headers: { "Access-Control-Allow-Origin": "*" },
      })
      .then((e) => {
        let arr = [];
        for (const [key, value] of Object.entries(e.data[0])) {
          if (value == "") {
            continue;
          }
          arr.push(value);
        }
        setWords(arr);
      });
  }, [day]);

  return (
    <div className="h-full">
      <div className="text-center">
        <div className="mt-5 font-semibold">Trending Topics</div>
      </div>
      <div className="divider"></div>
      <table className="table overflow-scroll max-h-5">
        <thead>
          <tr>
            <td>#</td>
            <td>Topic</td>
          </tr>
        </thead>
        <tbody>
          {words.map((e, i) => {
            return (
              <tr className="hover">
                <th>{i + 1}</th>
                <td onClick={() => history("/Search/" + e)}>{e}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default TendingTopcis;
