package models

// Post Struct: The main structure of the Json data. This contains all the information about an individual post.
type Post struct {
	ID               string       `json:"id"`
	CategoryID       string       `json:"categoryID"`
	Author           Author       `json:"author"`
	Title            string       `json:"title"`
	Body             string       `json:"body"`
	Anonymous        bool         `json:"anonymous"`
	Published        bool         `json:"published"`
	PublishedAt      string       `json:"publishedAt"`
	Group            string       `json:"group"`
	Number           int          `json:"number"`
	Type             string       `json:"type"`
	Visibility       string       `json:"visibility"`
	Slug             string       `json:"slug"`
	CreatedAt        string       `json:"createdAt"`
	UpdatedAt        string       `json:"updatedAt"`
	AnswersCount     int          `json:"answersCount"`
	UniqueViewsCount int          `json:"uniqueViewsCount"`
	ViewsCount       int          `json:"viewsCount"`
	AnsweredAt       string       `json:"AnsweredAt"`
	ModAnsweredAt    string       `json:"modAnsweredAt"`
	Read             bool         `json:"read"`
	Conversation     Conversation `json:"conversation"`
	Comments         []Comment    `json:"comments"`
}

// Author Struct: Contains information about an author- this can be the author of a post or a comment
type Author struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Registered bool   `json:"registered"`
	Slug       string `json:"slug"`
	Role       string `json:"role"`
}

// Conversation Struct: Contains information about the conversation associated with a post
type Conversation struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Slug          string       `json:"slug"`
	Photo         string       `json:"photo"`
	Type          string       `json:"type"`
	Public        bool         `json:"public"`
	Group         string       `json:"group"`
	Network       string       `json:"network"`
	FirstMessage  FirstMessage `json:"firstMessage"`
	LastMessage   LastMessage  `json:"lastMessage"`
	LastMessageAt string       `json:"lastMessageAt"`
	CreatedAt     string       `json:"createdAt"`
	UpdatedAt     string       `json:"updatedAt"`
	Post          string       `json:"post"`
	Subscribers   []Subscriber `json:"subscribers"`
}

// First Message Struct: A part of a conversation json. Data on the first message of the conversation
type FirstMessage struct {
	ID                string `json:"id"`
	UUID              string `json:"uuid"`
	AnonymousLevel    int    `json:"anonymousLevel"`
	Conversation      string `json:"conversation"`
	System            bool   `json:"system"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"UpdatedAt"`
	SystemMessageType string `json:"systemMessageType"`
	Type              string `json:"string"`
	Metadata          any    `json:"metadata"`
	ReadState         any    `json:"readState"`
}

// Last Message Struct: A part of the conversation json. Data on the last message of the conversation
type LastMessage struct {
	ID                string `json:"id"`
	UUID              string `json:"uuid"`
	AnonymousLevel    int    `json:"anonymousLevel"`
	Author            Author `json:"author"`
	Conversation      string `json:"conversation"`
	System            bool   `json:"system"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"UpdatedAt"`
	SystemMessageType string `json:"systemMessageType"`
	Type              string `json:"string"`
	Metadata          any    `json:"metadata"`
	ReadState         any    `json:"readState"`
}

// Subscriber Struct: Information of a subscriber, a part of the conversation struct, which contains and array of subscribers.
type Subscriber struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// Comment Struct: Information related to a comment. The Post json contains an array of comments
type Comment struct {
	ID          string `json:"id"`
	Author      Author `json:"author"`
	Body        string `json:"body"`
	Answer      bool   `json:"answer"`
	Metadata    any    `json:"metadata"`
	CreatedAt   string `json:"createdAt"`
	PublishedAt string `json:"publishedAt"`
	Endorsed    bool   `json:"endorsed"`
	Depth       int    `json:"depth"`
}

// UserData Struct: Information and statistics about a user. Add more fields when needed.
//
// For use, specific fields can be omitted if they are not needed.
type UserData struct {
	Author           Author `json:"author"`
	PostCount        int    `json:"postCount"`
	CommentCount     int    `json:"commentCount"`
	LastPostTime     string `json:"lastPostTime"`
	LastCommentTime  string `json:"lastCommentTime"`
	Role             string `json:"role"`
	EndorsedComments int    `json:"endorsedComments"`
}

// VerbatimUserData Struct: Composition of the UserData struct with vertabim data on posts and comments.
//
// Composition is essentially class inheritance.
type VerbatimUserData struct {
	UserData
	Posts    []Post    `json:"posts"`
	Comments []Comment `json:"comments"`
}

// CompactUserData Struct: Composition of the UserData struct with compact data on posts and comments.
//
// Composition is essentially class inheritance.
type CompactUserData struct {
	UserData
	Posts    []CompactPost    `json:"compactPosts"`
	Comments []CompactComment `json:"compactComments"`
}

type CompactPost struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Type        string `json:"type"`
	Visibility  string `json:"visibility"`
	Anonymous   bool   `json:"anonymous"`
	ModAnswered bool   `json:"modanswered"`
	ViewsCount  int    `json:"viewsCount"`
	PublishedAt string `json:"publishedAt"`
}

type CompactComment struct {
	ID          string `json:"id"`
	Body        string `json:"body"`
	Endorsed    bool   `json:"endorsed"`
	IsAnswer    bool   `json:"isAnswer"`
	PublishedAt string `json:"publishedAt"`
}

type DateData struct {
	Date      string `json:"date"`
	PostCount int    `json:"postCount"`
}

type ViewData struct {
	Date      string `json:"date"`
	ViewCount int    `json:"viewCount"`
}

type ViewCountData struct {
	Date  []string `json:"date"`
	Views []int    `json:"views"`
}

type PostCountData struct {
	Date  []string `json:"date"`
	Posts []int    `json:"posts"`
}

type AnswerData struct {
	AnsweredPosts   int `json:"AnsweredPosts"`
	UnansweredPosts int `json:"UnansweredPosts"`
}

type UnreadPosts struct {
	DateData    DateData `json:"dateData"`
	UnreadPosts []Post   `json:"UnreadPosts"`
}

type TopicWords struct {
	Word1        string `json:"word1"`
	Percentage1  string `json:"percentage1"`
	Word2        string `json:"word2"`
	Percentage2  string `json:"percentage2"`
	Word3        string `json:"word3"`
	Percentage3  string `json:"percentage3"`
	Word4        string `json:"word4"`
	Percentage4  string `json:"percentage4"`
	Word5        string `json:"word5"`
	Percentage5  string `json:"percentage5"`
	Word6        string `json:"word6"`
	Percentage6  string `json:"percentage6"`
	Word7        string `json:"word7"`
	Percentage7  string `json:"percentage7"`
	Word8        string `json:"word8"`
	Percentage8  string `json:"percentage8"`
	Word9        string `json:"word9"`
	Percentage9  string `json:"percentage9"`
	Word10       string `json:"word10"`
	Percentage10 string `json:"percentage10"`
}

type TopicWordsWithPercentages struct {
	Word1        string  `json:"word1"`
	Percentage1  float64 `json:"percentage1"`
	Word2        string  `json:"word2"`
	Percentage2  float64 `json:"percentage2"`
	Word3        string  `json:"word3"`
	Percentage3  float64 `json:"percentage3"`
	Word4        string  `json:"word4"`
	Percentage4  float64 `json:"percentage4"`
	Word5        string  `json:"word5"`
	Percentage5  float64 `json:"percentage5"`
	Word6        string  `json:"word6"`
	Percentage6  float64 `json:"percentage6"`
	Word7        string  `json:"word7"`
	Percentage7  float64 `json:"percentage7"`
	Word8        string  `json:"word8"`
	Percentage8  float64 `json:"percentage8"`
	Word9        string  `json:"word9"`
	Percentage9  float64 `json:"percentage9"`
	Word10       string  `json:"word10"`
	Percentage10 float64 `json:"percentage10"`
}

type TopicWordsWithoutPercentages struct {
	Word1  string `json:"word1"`
	Word2  string `json:"word2"`
	Word3  string `json:"word3"`
	Word4  string `json:"word4"`
	Word5  string `json:"word5"`
	Word6  string `json:"word6"`
	Word7  string `json:"word7"`
	Word8  string `json:"word8"`
	Word9  string `json:"word9"`
	Word10 string `json:"word10"`
}
