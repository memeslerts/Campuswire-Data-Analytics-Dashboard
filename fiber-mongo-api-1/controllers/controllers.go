package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"fiber-mongo-api/models"
	"fiber-mongo-api/mongodb"
)

var mongoClient = mongodb.MongoClient

// var maxTime = time.Unix(-2208988800, 0).Add(1<<63-1)
// var maxTimeFormat = maxTime.Format(time.RFC3339Nano)
// var validate = validator.New()

func GetPostById(c *fiber.Ctx) error {
	// Get api parameters
	postId := c.Params("postId")
	day := c.Params("day")

	// Query DB for post with the postId
	result, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": postId})
	}

	for _, post := range result {
		if post.ID == postId {
			return c.Status(http.StatusOK).JSON(post)
		}
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"Failed": "No Post ID"})
}

func GetPostsByDay(c *fiber.Ctx) error {
	// Get api parameters
	day := c.Params("day")

	// ?criteria=<string>&limit=<int>&reversed=<bool>
	criteria := c.Query("criteria", "views")
	limit := c.QueryInt("limit", 1)
	reversed := c.QueryBool("reversed")

	// Query DB for all posts up to a given day
	result, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	sortedResult, err := sortPostData(result, criteria, limit, reversed)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if limit > len(sortedResult) {
		return c.Status(http.StatusOK).JSON(sortedResult)
	}

	return c.Status(http.StatusOK).JSON(sortedResult[0:limit])
}

// sortUserData returns a list of user IDs sorted by the given criteria.
//
// The list of user IDs should be used with the map to get the UserData struct.
func sortPostData(postData []models.Post, criteria string, limit int, reversed bool) ([]models.Post, error) {
	// Map criteria to UserData field names
	criteriaMap := map[string]string{
		"views":         "ViewsCount",
		"uniqueViews":   "UniqueViewsCount",
		"answers":       "AnswersCount",
		"publishedAt":   "PublishedAt",
		"answeredAt":    "AnsweredAt",
		"modAnsweredAt": "ModAnsweredAt",
		"comments":      "Comments",
		"titleLength":   "Title",
		"bodyLength":    "Body",
	}

	criteriaField, exists := criteriaMap[criteria]

	if !exists {
		return nil, errors.New("invalid criteria")
	}

	filteredData := make([]models.Post, 0)
	for _, post := range postData {
		reflectPost := reflect.ValueOf(post)
		field := reflect.Indirect(reflectPost).FieldByName(criteriaField)

		switch criteria {
		case "views", "uniqueViews", "answers":
			if field.Int() > 0 {
				filteredData = append(filteredData, post)
			}
		case "publishedAt", "answeredAt", "modAnsweredAt":
			if field.String() != "" {
				filteredData = append(filteredData, post)
			}
		case "comments":
			if len(post.Comments) > 0 {
				filteredData = append(filteredData, post)
			}
		case "titleLength", "bodyLength":
			if len(field.String()) > 0 {
				filteredData = append(filteredData, post)
			}
		default:
			filteredData = append(filteredData, post)
		}
	}

	// // Get list of keys in engagementData
	// keys := make([]string, 0, len(filteredData))
	// for k := range filteredData {
	// 	keys = append(keys, k)
	// }

	// Sort keys by criteria
	sort.Slice(filteredData, func(i, j int) bool {
		postA := reflect.ValueOf(filteredData[i])
		postB := reflect.ValueOf(filteredData[j])
		fieldA := reflect.Indirect(postA).FieldByName(criteriaField)
		fieldB := reflect.Indirect(postB).FieldByName(criteriaField)

		var less bool
		switch criteria {
		case "views", "uniqueViews", "answers":
			less = fieldA.Int() > fieldB.Int()
		case "publishedAt", "answeredAt", "modAnsweredAt":
			less = fieldA.String() > fieldB.String()
		case "comments":
			less = fieldA.Len() > fieldB.Len()
		case "titleLength", "bodyLength":
			less = len(fieldA.String()) > len(fieldB.String())
		default:
			less = false
		}

		if reversed {
			return !less
		}
		return less
	})

	return filteredData, nil
}

func GetUserById(c *fiber.Ctx) error {
	day := c.Params("day")
	userId := c.Params("userid")

	// ?verbatim=<bool>
	verbatim := c.QueryBool("verbatim", false)

	// Query DB for all posts up to a given day
	result, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	// Abusing string comparison since [0-9] are always before [a-zA-Z]
	maxTimeString := "N/A"

	userData := models.UserData{
		Author:          models.Author{},
		PostCount:       0,
		CommentCount:    0,
		LastPostTime:    maxTimeString,
		LastCommentTime: maxTimeString,
		Role:            "student",
	}

	// Do NOT look at my bandaid implementation
	// *** TODO: Refractor to interfaces and dynamic dispatch ***
	if verbatim {
		user := models.VerbatimUserData{
			UserData: userData,
			Posts:    make([]models.Post, 0),
			Comments: make([]models.Comment, 0),
		}

		for _, post := range result {
			if post.Author.ID == userId {
				user.Author = post.Author
				user.Posts = append(user.Posts, post)
				user.PostCount += 1

				// Update last post time
				if post.PublishedAt < user.LastPostTime {
					user.LastPostTime = post.PublishedAt
				}
			}
			for _, comment := range post.Comments {
				// Update comment count
				if comment.Author.ID == userId {
					user.Comments = append(user.Comments, comment)
					user.CommentCount += 1

					if comment.Endorsed {
						user.EndorsedComments += 1
					}

					// Update last comment time
					if comment.PublishedAt < user.LastCommentTime {
						user.LastCommentTime = comment.PublishedAt
					}

					if isStaffComment(post, comment) {
						user.Role = "staff"
					}
				}
			}
		}
		return c.Status(http.StatusOK).JSON(user)
	} else {
		user := models.CompactUserData{
			UserData: userData,
			Posts:    make([]models.CompactPost, 0),
			Comments: make([]models.CompactComment, 0),
		}
		for _, post := range result {
			if post.Author.ID == userId {
				user.Author = post.Author
				compactPost := models.CompactPost{
					ID:          post.ID,
					Title:       post.Title,
					Body:        post.Body,
					Type:        post.Type,
					Visibility:  post.Visibility,
					Anonymous:   post.Anonymous,
					ModAnswered: post.ModAnsweredAt != "",
					ViewsCount:  post.ViewsCount,
					PublishedAt: post.PublishedAt,
				}
				user.Posts = append(user.Posts, compactPost)
				user.PostCount += 1

				// Update last post time
				if post.PublishedAt < user.LastPostTime {
					user.LastPostTime = post.PublishedAt
				}
			}
			for _, comment := range post.Comments {
				// Update comment count
				if comment.Author.ID == userId {
					compactComment := models.CompactComment{
						ID:          comment.ID,
						Body:        comment.Body,
						Endorsed:    comment.Endorsed,
						IsAnswer:    comment.Answer,
						PublishedAt: comment.PublishedAt,
					}
					user.Comments = append(user.Comments, compactComment)
					user.CommentCount += 1

					if comment.Endorsed {
						user.EndorsedComments += 1
					}

					// Update last comment time
					if comment.PublishedAt < user.LastCommentTime {
						user.LastCommentTime = comment.PublishedAt
					}

					if isStaffComment(post, comment) {
						user.Role = "staff"
					}
				}
			}
		}
		return c.Status(http.StatusOK).JSON(user)
	}
}

func GetEngagement(c *fiber.Ctx) error {
	day := c.Params("day")

	// ?criteria=<string>&limit=<int>&reversed=<bool>
	criteria := c.Query("criteria", "posts")
	limit := c.QueryInt("limit", 1)
	reversed := c.QueryBool("reversed")

	// Query DB for all posts up to a given day
	result, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	engagementData := getUserData(result)

	keys, err := sortUserData(engagementData, criteria, reversed)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	// Get a ranking of users
	ranking := make([]models.UserData, 0, limit)
	for i, k := range keys {
		ranking = append(ranking, engagementData[k])
		if i+1 >= int(limit) {
			break
		}
	}

	return c.Status(http.StatusOK).JSON(ranking)
}

func GetGraphPostCounts(c *fiber.Ctx) error {

	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])
	if err != nil {
		panic(err)
	}

	rng := c.QueryInt("range", 7)

	posts, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	dates := getDatesFromDays(dayAsInt, rng)
	postCountData := make(map[string]models.DateData)

	for _, post := range posts {
		publishedDate := post.PublishedAt[0:10]

		if slices.Contains(dates, publishedDate) {
			dateData, exists := postCountData[publishedDate]

			if !exists {
				dateData = models.DateData{Date: publishedDate, PostCount: 0}
			}

			dateData.PostCount += 1

			postCountData[publishedDate] = dateData
		}
	}

	var postDates []string
	var postCounts []int
	for _, date := range dates {
		postDates = append(postDates, date[5:])
		dateData, exists := postCountData[date]
		if !exists {
			postCounts = append(postCounts, 0)
		} else {
			postCounts = append(postCounts, dateData.PostCount)
		}
	}

	return c.Status(http.StatusOK).JSON(models.PostCountData{Date: postDates, Posts: postCounts})
}

func GetUniquePosts(c *fiber.Ctx) error {

	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])
	if err != nil {
		panic(err)
	}

	posts, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	date := getDatesFromDays(dayAsInt, 1)[0]
	var uniquePosts []models.Post

	for _, post := range posts {
		publishedDate := post.PublishedAt[0:10]

		if date == publishedDate {
			uniquePosts = append(uniquePosts, post)
		}
	}
	return c.Status(http.StatusOK).JSON(uniquePosts)
}
func GetRecentPosts(c *fiber.Ctx) error {
	// Get api parameters
	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])

	if err != nil {
		panic(err)
	}

	// Get an array of dates of previous 7 days from the queried date
	dates := getDatesFromDays(dayAsInt, 7)

	var title []string
	var authorName []string
	var viewCount []int
	var answerCount []int

	for i := 0; i < 7; i++ {

		if dayAsInt-(6-i) < 1 {
			continue
		}

		posts, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", "day"+strconv.Itoa(dayAsInt-(6-i)))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		for _, post := range posts {
			publishedDate := post.PublishedAt[0:10]

			if publishedDate == dates[i] {
				title = append(title, post.Title)
				authorName = append(authorName, (post.Author.FirstName + " " + post.Author.LastName))
				viewCount = append(viewCount, post.ViewsCount)
				answerCount = append(answerCount, post.AnswersCount)

			}
		}
	}

	// Returns an object with title, author, view count and answer count of all the posts in range
	type recentPosts struct {
		Names       []string
		AuthorNames []string
		Views       []int
		Answers     []int
	}
	return c.Status(http.StatusOK).JSON(recentPosts{Names: title, AuthorNames: authorName, Views: viewCount, Answers: answerCount})
}

func GetTopics(c *fiber.Ctx) error {
	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])

	topics, err := mongoClient.QueryAllML(bson.D{}, bson.D{}, 0, "campuswire-analytics", "updatedML"+strconv.Itoa(dayAsInt))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(topics)
}

func GetGraphViewCounts(c *fiber.Ctx) error {
	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])
	if err != nil {
		panic(err)
	}

	unique := c.QueryBool("unique", false)

	rng := c.QueryInt("range", 7)

	posts, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	dates := getDatesFromDays(dayAsInt, rng)
	viewCountData := make(map[string]models.ViewData)

	for _, post := range posts {
		publishedDate := post.PublishedAt[0:10]

		if slices.Contains(dates, publishedDate) {
			dateData, exists := viewCountData[publishedDate]

			if !exists {
				dateData = models.ViewData{Date: publishedDate, ViewCount: 0}
			}

			if unique {
				dateData.ViewCount += post.UniqueViewsCount
			} else {
				dateData.ViewCount += post.ViewsCount
			}

			viewCountData[publishedDate] = dateData
		}
	}

	var postDates []string
	var viewCounts []int
	for _, date := range dates {
		postDates = append(postDates, date[5:])
		dateData, exists := viewCountData[date]
		if !exists {
			viewCounts = append(viewCounts, 0)
		} else {
			viewCounts = append(viewCounts, dateData.ViewCount)
		}
	}

	return c.Status(http.StatusOK).JSON(models.ViewCountData{Date: postDates, Views: viewCounts})
}

func getDatesFromDays(currentDay int, rng int) []string {
	//returns an array of dates (strings in format yyyy-mm-dd) ending on day currentDay in a (inclusive) range of rng days
	//please ensure that this function is only called with days ranging between 1 and 92

	if currentDay < 1 || currentDay > 92 {
		panic("incorrect day")
	}

	var dates []string
	period := rng - 1

	if currentDay-period < -6 {
		period = currentDay + 6
	}
	for i := currentDay - period; i <= currentDay; i++ {
		dayDate := getDateFromDay(i)
		dates = append(dates, dayDate)
	}
	return dates
}

func getDateFromDay(day int) string {
	//returns an date (strings in format yyyy-mm-dd) corresponding to the specified day (day 1 contains 9-8 to 9-15, so we have days -6 to 0 as well)

	if day < -6 || day > 92 {
		panic("incorrect day range")
	}

	currentDay := time.Date(2022, 9, 8, 0, 0, 0, 0, time.UTC)
	currentDay = currentDay.Add(time.Hour * time.Duration(24*(day+6)))

	return currentDay.String()[0:10]
}

// getUserData returns a map of user IDs to UserData structs.
//
// It performs calculations based on the list of posts provided.
func getUserData(posts []models.Post) map[string]models.UserData {
	engagementData := make(map[string]models.UserData)

	// Abusing string comparison since [0-9] are always before [a-zA-Z]
	maxTimeString := "N/A"

	for _, post := range posts {
		authorData, exists := engagementData[post.Author.ID]

		if !exists {
			authorData = models.UserData{
				Author:           post.Author,
				PostCount:        0,
				CommentCount:     0,
				LastPostTime:     maxTimeString,
				LastCommentTime:  maxTimeString,
				Role:             "student",
				EndorsedComments: 0,
			}
		}

		// Update post count
		authorData.PostCount += 1

		// Update last post time
		if post.PublishedAt < authorData.LastPostTime {
			authorData.LastPostTime = post.PublishedAt
		}

		for _, comment := range post.Comments {
			// If comment is by the post author
			if comment.Author.ID == post.Author.ID {
				authorData.CommentCount += 1

				if comment.Endorsed {
					authorData.EndorsedComments += 1
				}

				// Update last comment time
				if comment.PublishedAt < authorData.LastCommentTime {
					authorData.LastCommentTime = comment.PublishedAt
				}

				if isStaffComment(post, comment) {
					authorData.Role = "staff"
				}

			} else { // If comment is by another user

				// Get comment author data
				commentAuthor, exists := engagementData[comment.Author.ID]

				if !exists {
					commentAuthor = models.UserData{
						Author:           comment.Author,
						PostCount:        0,
						CommentCount:     0,
						LastPostTime:     maxTimeString,
						LastCommentTime:  maxTimeString,
						Role:             "student",
						EndorsedComments: 0,
					}
				}

				commentAuthor.CommentCount += 1

				if comment.Endorsed {
					commentAuthor.EndorsedComments += 1
				}

				// Update last comment time
				if comment.PublishedAt < commentAuthor.LastCommentTime {
					commentAuthor.LastCommentTime = comment.PublishedAt
				}

				if isStaffComment(post, comment) {
					commentAuthor.Role = "staff"
				}

				engagementData[comment.Author.ID] = commentAuthor
			}
		}

		engagementData[post.Author.ID] = authorData
	}

	return engagementData
}

// isStaffComment returns true if the comment was posted by a staff member.
func isStaffComment(post models.Post, comment models.Comment) bool {
	return strings.Split(comment.PublishedAt, ".")[0] == strings.Split(post.ModAnsweredAt, ".")[0]
}

// sortUserData returns a list of user IDs sorted by the given criteria.
//
// The list of user IDs should be used with the map to get the UserData struct.
func sortUserData(userData map[string]models.UserData, criteria string, reversed bool) ([]string, error) {
	// Map criteria to UserData field names
	criteriaMap := map[string]string{
		"posts":            "PostCount",
		"comments":         "CommentCount",
		"postTime":         "LastPostTime",
		"commentTime":      "LastCommentTime",
		"endorsedComments": "EndorsedComments",
	}

	criteriaField, exists := criteriaMap[criteria]

	if !exists {
		return nil, errors.New("invalid criteria")
	}

	filteredData := make(map[string]models.UserData)
	for user := range userData {
		author := reflect.ValueOf(userData[user])
		field := reflect.Indirect(author).FieldByName(criteriaField)

		switch criteria {
		case "posts", "comments", "endorsedComments":
			if field.Int() > 0 {
				filteredData[user] = userData[user]
			}
		case "postTime", "commentTime":
			if field.String() != "N/A" {
				filteredData[user] = userData[user]
			}
		default:
			filteredData[user] = userData[user]
		}
	}

	// Get list of keys in engagementData
	keys := make([]string, 0, len(filteredData))
	for k := range filteredData {
		keys = append(keys, k)
	}

	// Sort keys by criteria
	sort.Slice(keys, func(i, j int) bool {
		authorA := reflect.ValueOf(filteredData[keys[i]])
		authorB := reflect.ValueOf(filteredData[keys[j]])
		fieldA := reflect.Indirect(authorA).FieldByName(criteriaField)
		fieldB := reflect.Indirect(authorB).FieldByName(criteriaField)

		var less bool
		switch criteria {
		case "posts", "comments", "endorsedComments":
			less = fieldA.Int() > fieldB.Int()
		case "postTime", "commentTime":
			less = fieldA.String() > fieldB.String()
		default:
			less = false
		}

		if reversed {
			return !less
		}
		return less
	})

	return keys, nil
}
func GetAnsweredPosts(c *fiber.Ctx) error {
	// Get api parameters
	day := c.Params("day")

	// Query DB for all posts up to a given day
	results, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	var questionsData models.AnswerData

	for _, post := range results {
		if post.AnswersCount > 0 {
			questionsData.AnsweredPosts += 1
		} else {
			questionsData.UnansweredPosts += 1
		}
	}

	return c.Status(http.StatusOK).JSON(questionsData)
}

func GetUnreadPosts(c *fiber.Ctx) error {
	// Get api parameters
	day := c.Params("day")
	dayAsInt, err := strconv.Atoi(day[3:])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	qrange := c.QueryInt("range", 7)

	// Query DB for all posts up to a given day
	results, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	dates := getDatesFromDays(dayAsInt, qrange)

	unreadPostsData := make(map[string]models.UnreadPosts)

	for _, post := range results {
		publishedDate := post.PublishedAt[0:10]
		if slices.Contains(dates, publishedDate) && !post.Read && post.AnswersCount < 1 {
			unreadPost := unreadPostsData[publishedDate]
			unreadPost.DateData.PostCount += 1
			unreadPost.UnreadPosts = append(unreadPost.UnreadPosts, post)
			unreadPostsData[publishedDate] = unreadPost
		}
	}

	return c.Status(http.StatusOK).JSON(unreadPostsData)
}

func SearchForToken(c *fiber.Ctx) error {
	// Get api parameters
	day := c.Params("day")
	searchToken := c.Query("token", "")
	deleteCost := c.QueryInt("deleteCost", 1)
	insertCost := c.QueryInt("insertCost", 1)
	replaceCost := c.QueryInt("replaceCost", 1)
	similarityThreshold := c.QueryFloat("similarityThreshold", 0.6)
	similarityDistance := c.QueryInt("similarityDistance", 3)
	similarityMethod := c.Query("similarityMethod", "threshold")

	if similarityMethod != "distance" && similarityMethod != "threshold" {
		similarityMethod = "threshold"
	}

	// Query DB for all posts up to a given day
	result, err := mongoClient.QueryAll(bson.D{}, bson.D{}, 0, "campuswire-analytics", day)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	type SearchItems struct {
		Posts    []models.Post
		Comments []models.Comment
		Users    []models.Author
	}
	type SearchResults struct {
		ExactSearchItems   SearchItems
		SimilarSearchItems SearchItems
	}

	searchResults := SearchResults{
		ExactSearchItems: SearchItems{
			Posts:    make([]models.Post, 0),
			Comments: make([]models.Comment, 0),
			Users:    make([]models.Author, 0),
		},
		SimilarSearchItems: SearchItems{
			Posts:    make([]models.Post, 0),
			Comments: make([]models.Comment, 0),
			Users:    make([]models.Author, 0),
		},
	}
	similarUsers := make(map[string]models.Author)
	exactUsers := make(map[string]models.Author)

	searchToken = strings.ToLower(searchToken)

	exactSearchRegex := regexp.MustCompile(`"[^"]*"`)
	exactSearchTokens := make([]string, 0)

	exactSearchString := exactSearchRegex.FindString(searchToken)
	for exactSearchString != "" {
		exactSearchTokens = append(exactSearchTokens, strings.Trim(exactSearchString, `"`))
		before, after, found := strings.Cut(searchToken, exactSearchString)
		if found {
			searchToken = before + after
		}

		exactSearchString = exactSearchRegex.FindString(searchToken)
	}

	searchWords := strings.Fields(searchToken)

	lev := metrics.NewLevenshtein()
	lev.CaseSensitive = false
	lev.DeleteCost = deleteCost
	lev.InsertCost = insertCost
	lev.ReplaceCost = replaceCost

	isSimilarToken := func(token string, searchToken string) bool {
		if similarityMethod == "distance" {
			return lev.Distance(token, searchToken) <= similarityDistance
		} else {
			return strutil.Similarity(token, searchToken, lev) > similarityThreshold
		}
	}

	// Useful for debugging purposes
	similarTokens := make([]string, 0)
	exactTokens := make([]string, 0)

	for _, post := range result {
		postSearchText := strings.ToLower(post.Title + " " + post.Body)
		nameSearchText := strings.ToLower(post.Author.FirstName + " " + post.Author.LastName)

		// Exact Search
		for _, exactSearchToken := range exactSearchTokens {
			if strings.Contains(postSearchText, exactSearchToken) {
				exactTokens = append(exactTokens, postSearchText)
				searchResults.ExactSearchItems.Posts = append(searchResults.ExactSearchItems.Posts, post)
			}
			if strings.Contains(nameSearchText, exactSearchToken) {
				exactTokens = append(exactTokens, nameSearchText)
				exactUsers[post.Author.ID] = post.Author
			}
		}

		// Similar Search
		postTokens := strings.Fields(strings.ToLower(postSearchText))
		nameTokens := strings.Fields(strings.ToLower(nameSearchText))
		for _, searchWord := range searchWords {
			for _, token := range postTokens {
				if isSimilarToken(token, searchWord) {
					similarTokens = append(similarTokens, token)
					searchResults.SimilarSearchItems.Posts = append(searchResults.SimilarSearchItems.Posts, post)
				}
			}
			for _, token := range nameTokens {
				if isSimilarToken(token, searchWord) {
					similarTokens = append(similarTokens, token)
					similarUsers[post.Author.ID] = post.Author
				}
			}
		}
		for _, comment := range post.Comments {
			commentSearchText := strings.ToLower(comment.Body)
			commentNameSearchText := strings.ToLower(comment.Author.FirstName + " " + comment.Author.LastName)

			// Exact Search
			for _, exactSearchToken := range exactSearchTokens {
				if strings.Contains(commentSearchText, exactSearchToken) {
					exactTokens = append(exactTokens, commentSearchText)
					searchResults.ExactSearchItems.Comments = append(searchResults.ExactSearchItems.Comments, comment)
				}
				if strings.Contains(commentNameSearchText, exactSearchToken) {
					exactTokens = append(exactTokens, commentNameSearchText)
					exactUsers[comment.Author.ID] = comment.Author
				}
			}

			// Similar Search
			commentTokens := strings.Fields(strings.ToLower(commentSearchText))
			commentNameTokens := strings.Fields(strings.ToLower(commentNameSearchText))
			for _, searchWord := range searchWords {
				for _, token := range commentTokens {
					if isSimilarToken(token, searchWord) {
						similarTokens = append(similarTokens, token)
						searchResults.SimilarSearchItems.Comments = append(searchResults.SimilarSearchItems.Comments, comment)
					}
				}
				for _, token := range commentNameTokens {
					if isSimilarToken(token, searchWord) {
						similarTokens = append(similarTokens, token)
						similarUsers[comment.Author.ID] = comment.Author
					}
				}
			}
		}
	}

	// Convert to lists
	for _, author := range similarUsers {
		searchResults.SimilarSearchItems.Users = append(searchResults.SimilarSearchItems.Users, author)
	}
	for _, author := range exactUsers {
		searchResults.ExactSearchItems.Users = append(searchResults.ExactSearchItems.Users, author)
	}

	return c.Status(http.StatusOK).JSON(searchResults)
}
