import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

function Post() {
  const navigate = useNavigate();
  const { postId } = useParams();
  const [post, loadPost] = useState({
    author: {id: "", firstName: "", lastName: "", registered: false, slug: "", role: ""},
    title: "",
    body: "",
    publishedAt: "",
    viewsCount: 0,
    comments: [],
  });
  useEffect(() => {
    axios
      .get("http://localhost:8080/day92/post/" + postId, {
        headers: { "Access-Control-Allow-Origin": "*" },
      })
      .then((e) => {
        loadPost(e.data);
      });
  }, []);

  return (
    <div>
      <div className="flex text-right">
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
          {post.author.firstName + " " + post.author.lastName}
        </div>
        <div className="text-center italic">Published At: {post.publishedAt}</div>
        <div className="text-center italic">Views: {post.viewsCount}</div>
        <br></br><br></br>
        <div className="text-center text-4xl">{post.title}</div>
        <br></br>
        <div className="px-10">{post.body}</div>
      </div>
      <div className="flex flex-wrap justify-around">
        {post.comments.map((comment) => {
          return (
            <div className="card w-96 bg-base-100 shadow-xl mt-10">
              <div className="card-body">
                <h2 className="card-title">{comment.author.firstName + comment.author.lastName}</h2>
                <p>{comment.body}</p>
                <div className="text-sm text-gray-400 text-end">
                  Endorsed: {String(comment.endorsed)}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default Post;
