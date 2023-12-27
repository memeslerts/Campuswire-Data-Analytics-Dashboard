import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

function User() {
  const navigate = useNavigate();
  const [dropdown, showDropdown] = useState(true);
  const { userId } = useParams();
  const [user, loadUser] = useState({
    author: { firstName: "", lastName: "" },
    postCount: 0,
    commentCount: 0,
    lastCommentTime: "",
    lastPostTime: "",
    compactPosts: [],
  });
  let params = userId.split(",")
  const [day, selectDay] = useState(params[1]);
  useEffect(() => {
    axios
      .get("http://localhost:8080/day" + day + "/user/" + params[0], {
        headers: { "Access-Control-Allow-Origin": "*" },
      })
      .then((e) => {
        loadUser(e.data);
      });
  }, [day]);

  function listOfDays() {
    let days = [];
    for (let i = 1; i < 90; i++) {
      days.push(
        <li className="w-full">
          <button
            onClick={() => {
              selectDay(i);
              showDropdown(false);
              axios
                .get("http://localhost:8080/day" + i + "/user/" + params[0], {
                  headers: { "Access-Control-Allow-Origin": "*" },
                })
                .then((e) => {
                  loadUser(e.data);
                });
            }}
          >
            Day {i}
          </button>
        </li>
      );
    }
    return days;
  }

  return user.author.id == "" ? (
    <div className="font-bold text-5xl grid place-content-center w-screen h-screen">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth={1.5}
        stroke="currentColor"
        className="w-50 h-50"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
        />
      </svg>
      <div>User Not Found</div>
    </div>
  ) : (
    <div>
      <div className="flex">
        <div className="dropdown w-full grow">
          <div
            tabIndex="0"
            role="button"
            className="btn w-full rounded-none"
            onClick={() => showDropdown(true)}
          >
            Day 1 - Day {day}
          </div>
          {dropdown && (
            <ul className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52 max-h-52 overflow-x-scroll flex flex-row">
              {listOfDays()}
            </ul>
          )}
        </div>
        <button
          className="btn btn-square rounded-none"
          onClick={() => navigate("/")}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
      <div className="text-center">
        <div className="text-center text-3xl font-semibold mt-5">
          {user.author.firstName + " " + user.author.lastName}
        </div>
        <div className="text-center italic">Id: {userId}</div>
        <div className="badge badge-primary badge-outlinem mt-5">student</div>
        <div className="mt-5">
          <div className="stats shadow">
            <div className="stat">
              <div className="stat-figure text-secondary">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="inline-block w-8 h-8 stroke-current"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
              </div>
              <div className="stat-title">Post Count</div>
              <div className="stat-value">{user.postCount}</div>
              <div className="stat-desc">{user.lastPostTime}</div>
            </div>

            <div className="stat">
              <div className="stat-figure text-secondary">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="inline-block w-8 h-8 stroke-current"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                  ></path>
                </svg>
              </div>
              <div className="stat-title">Comment Count</div>
              <div className="stat-value">{user.commentCount}</div>
              <div className="stat-desc">{user.lastCommentTime}</div>
            </div>
          </div>
        </div>
      </div>
      <div className="flex flex-wrap justify-around">
        {user.compactPosts.map((post) => {
          return (
            <div className="card w-96 bg-base-100 shadow-xl mt-10">
              <div className="card-body">
                <h2 className="card-title">{post.title}</h2>
                <p>{post.body}</p>
                <div className="text-sm text-gray-400 text-end">
                  Views: {post.viewsCount}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default User;
