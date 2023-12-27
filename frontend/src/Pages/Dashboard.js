import { useNavigate } from "react-router-dom";
import EngagementTable from "../Components/EngagementTable";
import TrendingTopics from "../Components/TrendingTopics";
import RecentPosts from "../Components/RecentPosts";
import GraphSlider from "../Components/GraphSlider";
import Search from "../Components/Search";
import { useState } from "react";

function Dashboard() {
  const navigate = useNavigate();
  const [day, selectDay] = useState(1);
  const [dropdown, showDropdown] = useState(true);

  function listOfDays() {
    let days = [];
    for (let i = 1; i < 90; i++) {
      days.push(
        <li className="w-full">
          <button
            onClick={() => {
              selectDay(i);
              showDropdown(false);
            }}
          >
            Day {i}
          </button>
        </li>
      );
    }
    return days;
  }

  return (
    <div className="w-screen h-screen flex flex-col">
      <div className="dropdown w-full">
        <div
          tabIndex="0"
          role="button"
          className="btn w-full rounded-none rounded-t-lg"
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
      <div className="w-full h-full p-5 flex flex-col gap-10">
        <div className="flex gap-4 items-center">
          <div className="text-5xl font-bold text-blue-400">Welcome Back!</div>
          <div className="w-3/4 absolute left-1/2 -translate-x-1/2 z-[1]"><Search day={day} /></div>
        </div>
        <div className="h-full grid grid-cols-6 gap-5">
          <div className="col-span-1 w-full h-full outline outline-4 outline-gray-300 rounded-lg shadow-xl">
            <EngagementTable day={day} />
          </div>
          <div className="grid grid-col-2 col-span-4 grid-rows-2 gap-y-3 relative">
            <div className="col-span-1 outline outline-4 outline-gray-300 rounded-lg shadow-xl">
              <RecentPosts day={day} />
            </div>
            <div className="col-span-1 flex flex-row justify-around">
              <GraphSlider sliderType="Pie" day={day} />
              <GraphSlider sliderType="Bar" day={day} />
            </div>
          </div>
          <div className="col-span-1 w-full h-full outline outline-4 outline-gray-300 rounded-lg shadow-xl">
            <TrendingTopics day={day} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
