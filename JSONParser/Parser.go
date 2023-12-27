package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	//"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func main() {
	// Connecting to MongoDB. Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// IMPORTANT: IF YOU USE THIS LINE TO CONNECT TO DATABASE: CHANGE USERNAME AND PASSWORD WITH YOUR INFO
	opts := options.Client().ApplyURI("mongodb+srv://javin_m:MongoJav123@campuswire-analytics.wxftjuk.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Open directory of all data
	searchDir := "./Data"

	// Create a fileList from the directory
	fileList := make([]string, 0)
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})
	if e != nil {
		panic(e)
	}

	// counter of which file we are on
	day := 0

	for _, file := range fileList {

		fmt.Printf("%s\n", file)

		//Open a specific json file. Print any error that occurs in the process
		jsonFile, err := os.Open(file)

		if err != nil {
			fmt.Println(err)
			fmt.Println("")
		} else {
			fmt.Println("Successfully opened json!")
		}
		fmt.Println("")

		//Defer closing the json file until the end of main (makes the call jsonFile.Close() run at the end of main)
		defer jsonFile.Close()

		//Read the file and convert it into a byte array. Print any error if it occurs
		json_response, err := io.ReadAll(jsonFile)

		if err != nil {
			fmt.Printf("%s\r\n", err)
		}

		//Create a variable (space in memory) to store the parsed data. Stored in an array of the Post struct.
		var posts []Post

		//Unmarshall (unravel) the byte array into the memory address of the created variable for storing. Check and print any errors in the process
		err = json.Unmarshal(json_response, &posts)
		if err != nil {
			fmt.Printf("%s\r\n", err)
		}

		//A simple loop that loops through each post, printing some info related to it
		//Feel free to edit or change the printed values (look at the Structs to see what can be printed)

		/*
			for i := 0; i < len(posts); i++ {

				fmt.Println("----------------------------")

				fmt.Println("First Name: " + posts[i].Author.FirstName)
				fmt.Println("Last Name: " + posts[i].Author.LastName)
				fmt.Println("Anonymous: " + strconv.FormatBool(posts[i].Anonymous))

				//Convert the PublishedAt string (formatted as RFC3339) into a Time.time type. Check/print any errors in the process
				date, error := time.Parse(time.RFC3339, posts[i].PublishedAt)
				if error != nil {
					fmt.Println(error)
				} else {
					fmt.Printf("Posted At: %v\n", date)
				}

				fmt.Println("Title: " + posts[i].Title)
				fmt.Println("Body: " + posts[i].Body)

				fmt.Println("----------------------------")
			}
		*/

		// Specify the Collection in which to insert the data. Currently putting each file (day) into a new collection
		coll := client.Database("campuswire-analytics").Collection("day" + strconv.Itoa(day))

		for i := 0; i < len(posts); i++ {
			doc := posts[i]

			// Insert post into collection
			result, err := coll.InsertOne(context.TODO(), doc)
			if err != nil {
				fmt.Printf("%s\r\n%s\r\n", err, result)
			}

			//fmt.Printf("Inserted document with _idL %v\n", result.InsertedID)
		}

		day += 1
	}
}
