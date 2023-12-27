package routes

import (
	"github.com/gofiber/fiber/v2"

	"fiber-mongo-api/controllers"
)

func UserRoute(app *fiber.App) {
	// All routes come here
	day := app.Group("/:day")
	day.Get("post/:postId", controllers.GetPostById)
	day.Get("posts", controllers.GetPostsByDay)
	day.Get("engagement", controllers.GetEngagement)
	day.Get("posts/unique", controllers.GetUniquePosts)
	day.Get("recent", controllers.GetRecentPosts)
	day.Get("search", controllers.SearchForToken)

	// User Page
	user := day.Group("/user")
	user.Get("/:userid", controllers.GetUserById)

	// Dashboard features
	dashboard := day.Group("/dashboard")
	dashboard.Get("engagement", controllers.GetEngagement)
	dashboard.Get("topics", controllers.GetTopics)

	//Graph features
	graph := dashboard.Group("/graph")
	graph.Get("postCount", controllers.GetGraphPostCounts)
	graph.Get("viewCount", controllers.GetGraphViewCounts)
	graph.Get("answerData", controllers.GetAnsweredPosts)
	graph.Get("unread", controllers.GetUnreadPosts)

}
