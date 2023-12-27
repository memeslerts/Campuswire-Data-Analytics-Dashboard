import { BrowserRouter, Routes, Route } from "react-router-dom";
import GraphSlider from "./Components/GraphSlider";
import EngagementTable from "./Components/EngagementTable";
import RecentPosts from "./Components/RecentPosts";
import Dashboard from "./Pages/Dashboard";
import Post from "./Pages/Post";
import User from "./Pages/User";
import SearchPage from "./Pages/SearchPage";

function Router() {
    return ( 
        <BrowserRouter>
            <Routes>
                <Route path="/" Component={Dashboard} />
                <Route path="/GraphSlider" Component={GraphSlider} />
                <Route path="/EngagementTable" Component={EngagementTable} />
                <Route path="/RecentPosts" Component={RecentPosts} />
                <Route path="/Post">
                    <Route path=":postId" Component={Post} />
                </Route>
                <Route path="/User">
                    <Route path=":userId" Component={User}/>
                </Route>
		<Route path="/Search">
                    <Route path=":query" Component={SearchPage}/>
                </Route>
            </Routes>
        </BrowserRouter>
     );
}

export default Router;