package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fiber-mongo-api/models"
)

type Mongo struct {
	client *mongo.Client // raw mongo client
}

// Setup returns a general mongo wrapper by loading URI from environment variables.
func Setup() Mongo {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Get URI from environment variables
	MONGO_URI := os.Getenv("MONGO_URI")
	return New(MONGO_URI)
}

// New returns a new mongo wrapper object from a connection URI.
func New(URI string) Mongo {
	// Set mongo client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	mongo_client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Wrap raw client in Mongo struct
	return Mongo{mongo_client}
}

// QueryAll returns a list of documents that match the query.
//
// The collection is first sorted, then queried based on the given filter.
// The limit parameter is the maximum number of documents returned.
// To get all results, set limit to 0.
func (m Mongo) QueryAll(query bson.D, sort bson.D, limit int, databaseName string, collectionName string) ([]models.Post, error) {
	// Get database and collection from names
	database := m.client.Database(databaseName)
	collection := database.Collection(collectionName)

	// Allocate a list of Posts for the query
	var documents []models.Post

	// Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set limit and sort options
	find_opts := options.Find().SetLimit(int64(limit)).SetSort(sort)
	cursor, err := collection.Find(ctx, query, find_opts)
	if err != nil {
		return documents, err
	}

	// Iterate through cursor and decode results as models.Post
	for cursor.Next(ctx) {
		var result models.Post
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		result = cleanPostData(result)
		documents = append(documents, result)
	}
	return documents, err
}

func (m Mongo) QueryAllML(query bson.D, sort bson.D, limit int, databaseName string, collectionName string) ([]models.TopicWordsWithoutPercentages, error) {
	// Get database and collection from names
	database := m.client.Database(databaseName)
	collection := database.Collection(collectionName)

	// Allocate a list of Posts for the query
	var documents []models.TopicWordsWithoutPercentages

	// Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set limit and sort options
	find_opts := options.Find().SetLimit(int64(limit)).SetSort(sort)
	cursor, err := collection.Find(ctx, query, find_opts)
	if err != nil {
		return documents, err
	}

	collectionCount := 0
	// Iterate through cursor and decode results as models.TopicWords
	for cursor.Next(ctx) {
		collectionCount += 1
		var result models.TopicWords
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		documents = append(documents, cleanTopicData(result))
	}
	if collectionCount == 0 {
		return documents, errors.New("0 ITEMS IN COLLECTION")
	}
	return documents, err
}

func cleanTopicData(topic models.TopicWords) models.TopicWordsWithoutPercentages {
	var topicA models.TopicWordsWithoutPercentages
	//var newTopic models.TopicWordsWithPercentages
	var err error

	topicA.Word1 = cleanTopicString(topic.Word1, true)
	topicA.Word2 = cleanTopicString(topic.Word2, true)
	topicA.Word3 = cleanTopicString(topic.Word3, true)
	topicA.Word4 = cleanTopicString(topic.Word4, true)
	topicA.Word5 = cleanTopicString(topic.Word5, true)
	topicA.Word6 = cleanTopicString(topic.Word6, true)
	topicA.Word7 = cleanTopicString(topic.Word7, true)
	topicA.Word8 = cleanTopicString(topic.Word8, true)
	topicA.Word9 = cleanTopicString(topic.Word9, true)
	topicA.Word10 = cleanTopicString(topic.Word10, true)

	/*newTopic.Word1 = cleanTopicString(topic.Word1, false)
	if len(newTopic.Word1) > 0 {
		newTopic.Percentage1, err = strconv.ParseFloat(topic.Percentage1, 64)
	}
	newTopic.Word2 = cleanTopicString(topic.Word2, false)
	if len(newTopic.Word2) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage2[1:], 64)
	}
	newTopic.Word3 = cleanTopicString(topic.Word3, false)
	if len(newTopic.Word3) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage3[1:], 64)
	}
	newTopic.Word4 = cleanTopicString(topic.Word4, false)
	if len(newTopic.Word4) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage4[1:], 64)
	}
	newTopic.Word5 = cleanTopicString(topic.Word5, false)
	if len(newTopic.Word5) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage5[1:], 64)
	}
	newTopic.Word6 = cleanTopicString(topic.Word6, false)
	if len(newTopic.Word6) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage6[1:], 64)
	}
	newTopic.Word7 = cleanTopicString(topic.Word7, false)
	if len(newTopic.Word7) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage7[1:], 64)
	}
	newTopic.Word8 = cleanTopicString(topic.Word8, false)
	if len(newTopic.Word8) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage8[1:], 64)
	}
	newTopic.Word9 = cleanTopicString(topic.Word9, false)
	if len(newTopic.Word9) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage9[1:], 64)
	}
	newTopic.Word10 = cleanTopicString(topic.Word10, true)
	if len(newTopic.Word10) > 0 {
		newTopic.Percentage2, err = strconv.ParseFloat(topic.Percentage10[1:], 64)
	}*/

	if err != nil {
		fmt.Errorf(err.Error())
	}

	return topicA
}

func cleanTopicString(s string, last bool) string {
	if len(s) <= 2 {
		return ""
	} else if s[0] == '"' && s[len(s)-1] == ' ' && s[len(s)-2] == '"' {
		return s[1 : len(s)-2]
	}
	return s[1 : len(s)-1]
}

// QueryOne returns the first document that match the query.
//
// The collection is first sorted, then queried based on the given filter.
func (m Mongo) QueryOne(query bson.D, sort bson.D, databaseName string, collectionName string) (models.Post, error) {
	// Get database and collection from names
	database := m.client.Database(databaseName)
	collection := database.Collection(collectionName)

	// Allocate a Post struct for the query
	var document models.Post

	// Timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set sort option
	opts := options.FindOne().SetSort(sort)
	err := collection.FindOne(ctx, query, opts).Decode(&document)

	document = cleanPostData(document)

	return document, err
}

// RawClient returns the raw mongo-driver client.
//
// This is useful if you need to do something that is not supported by the wrapper.
func (m Mongo) RawClient() *mongo.Client {
	return m.client
}

// Close disconnects the client from the MongoDB server.
//
// This should be defered immediately to properly clean up resources.
func (m Mongo) Close() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func cleanPostText(text string) string {
	newText := text
	if strings.HasPrefix(text, "zzz ") {
		newText = text[4:]
	}
	return newText
}

func toDisplayableDate(date string) string {
	if date == "" {
		return date
	}
	tokens := strings.Split(date, "T")
	return tokens[0] + " " + strings.Split(tokens[1], ".")[0]
}

func cleanPostData(post models.Post) models.Post {
	post.Title = cleanPostText(post.Title)
	post.Body = cleanPostText(post.Body)
	post.PublishedAt = toDisplayableDate(post.PublishedAt)
	post.CreatedAt = toDisplayableDate(post.CreatedAt)
	post.UpdatedAt = toDisplayableDate(post.UpdatedAt)
	post.AnsweredAt = toDisplayableDate(post.AnsweredAt)
	post.ModAnsweredAt = toDisplayableDate(post.ModAnsweredAt)
	post.Conversation.LastMessageAt = toDisplayableDate(post.Conversation.LastMessageAt)
	post.Conversation.CreatedAt = toDisplayableDate(post.Conversation.CreatedAt)
	post.Conversation.UpdatedAt = toDisplayableDate(post.Conversation.UpdatedAt)
	post.Conversation.FirstMessage.CreatedAt = toDisplayableDate(post.Conversation.FirstMessage.CreatedAt)
	post.Conversation.FirstMessage.UpdatedAt = toDisplayableDate(post.Conversation.FirstMessage.UpdatedAt)
	post.Conversation.LastMessage.CreatedAt = toDisplayableDate(post.Conversation.LastMessage.CreatedAt)
	post.Conversation.LastMessage.UpdatedAt = toDisplayableDate(post.Conversation.LastMessage.UpdatedAt)

	cleanedComments := []models.Comment{}
	for _, comment := range post.Comments {
		comment.Body = cleanPostText(comment.Body)
		comment.CreatedAt = toDisplayableDate(comment.CreatedAt)
		comment.PublishedAt = toDisplayableDate(comment.PublishedAt)
		cleanedComments = append(cleanedComments, comment)
	}
	post.Comments = cleanedComments
	return post
}

// Mongo wrapper for general use.
// If multiple clients are necessary, use New() instead.
var MongoClient = Setup()
