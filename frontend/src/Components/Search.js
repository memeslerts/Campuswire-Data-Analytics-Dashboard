import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function Search({ day }) {
  const [search, setSearch] = useState([]);
  const [input_value, setInputValue] = useState("");
  const [showDropDown, setShowDropDown] = useState(false);
  const history = useNavigate();

  function startSearch() {
    axios
      .get("http://localhost:8080/day" + day + "/search?token=" + input_value)
      .then((e) => {
        let arr = [];
        let types_of_searches = ["ExactSearchItems", "SimilarSearchItems"];
        let count = { Users: 0, Posts: 0, Comments: 0 };
        for (const type of types_of_searches) {
          for (const [key, value] of Object.entries(e.data[type])) {
            for (const item of value) {
              if (count[key] > 5) {
                continue;
              }
              count[key] = count[key] + 1;
              if (key == "Posts") {
                arr.push({
                  id: item.id,
                  key: key,
                  first_name: item.author.firstName,
                  last_name: item.author.lastName,
                  title: item.title,
                });
              } else if (key == "Comments") {
                let short_body = item.body
                if (item.body.length > 60){
                  short_body = item.body.substring(0, 58) + '...'
                }
                arr.push({
                  id: item.id,
                  key: key,
                  first_name: item.author.firstName,
                  last_name: item.author.lastName,
                  body: short_body,
                });
              } else {
                arr.push({
                  id: item.id,
                  key: key,
                  first_name: item.firstName,
                  last_name: item.lastName,
                });
              }
            }
          }
        }
        setSearch(arr);
      });
  }

  return (
    <div className="grow text-center">
      <span className="dropdown w-2/4">
        <input
          onChange={(e) => setInputValue(e.target.value)}
          onBlur={() => setShowDropDown(false)}
          placeholder="Search"
          className="px-4 py-3 outline outline-2 outline-gray-300 focus:outline-blue-300 rounded-lg w-full shadow-md m-1"
        />
        {showDropDown && (
          <ul className="absolute z-[1] menu p-2 bg-base-100 shadow rounded-box w-full">
            {search.map((e) => {
              if (e.key == "Users") {
                return (
                  <li onClick={() => history("/User/" + e["id"] + "," + day)}>
                    <p>
                      <span className="font-semibold">{e.first_name} {e.last_name}</span>
                      <span className="text-end italic">{e.key}</span>
                    </p>
                  </li>
                );
              } else if (e.key == "Posts") {
                return (
                  <li onClick={() => history("/Post/" + e["id"])}>
                    <p>
                      <span className="font-semibold">{e.title}</span>
                      <span className="text-end italic">{e.key}</span>
                    </p>
                  </li>
                );
              } else if (e.key == "Comments") {
                return (
                  <li>
                    <p>
                      <span className="font-semibold">{e.body}</span>
                      <span className="text-end italic">{e.key}</span>
                    </p>
                  </li>
                );
              }
            })}
          </ul>
        )}
      </span>
      <button
        className="btn ml-4"
        onClick={(e) => {
          startSearch();
          setShowDropDown(true);
        }}
      >
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
  );
}

export default Search;
