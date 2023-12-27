import { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

function SearchPage({ day }) {
  const [result, setResult] = useState([
    {
      type: "",
      author: "",
      title: "",
      body: "",
    },
  ]);
  const [input_value, setInputValue] = useState("");
  const { query } = useParams();
  const history = useNavigate();

  useEffect(() => {
    axios
      .get("http://localhost:8080/day" + 90 + "/search?token=" + query)
      .then((e) => {
        let arr = [];
        let types_of_searches = ["ExactSearchItems", "SimilarSearchItems"];
        let count = { Users: 0, Posts: 0, Comments: 0 };
        for (const type of types_of_searches) {
          for (const [key, value] of Object.entries(e.data[type])) {
            for (const item of value) {
              let obj = {
                id: item.id,
                key: key,
              };
              if (key == "Users") {
                obj["first_name"] = item.firstName;
                obj["last_name"] = item.lastName;
              } else if (key == "Posts") {
                obj["first_name"] = item.author.firstName;
                obj["last_name"] = item.author.lastName;
                obj["title"] = item.title;
              } else {
                obj["first_name"] = item.author.firstName;
                obj["last_name"] = item.author.lastName;
                obj["body"] = item.body;
              }

              arr.push(obj);
            }
          }
        }
        setResult(arr);
      });
  }, [query]);

  return (
    <div className="mt-5">
      <div className="text-center">
        <input
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Search"
          className="px-4 py-3 outline outline-2 outline-gray-300 focus:outline-blue-300 w-1/2 rounded-lg shadow-md m-1"
        />
        <button className="btn ml-4" onClick={() => history("/Search/" + input_value)}>
          Enter{" "}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M9 15L3 9m0 0l6-6M3 9h12a6 6 0 010 12h-3"
            />
          </svg>
        </button>
      </div>
      <div className="overflow-x-auto mt-10">
        <table className="table table-xs">
          <thead>
            <tr>
              <th></th>
              <th>Type</th>
              <th>Author</th>
              <th>Content</th>
              <th>id</th>
            </tr>
          </thead>
          <tbody>
            {result.map((e, i) => {
              return (
                <tr>
                  <th>{i}</th>
                  <td>{e.key}</td>
                  <td>
                    {e.first_name} {e.last_name}
                  </td>
                  <td>{e.key == "Posts" ? e.title : e.body}</td>
                  <td>{e.id}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default SearchPage;
