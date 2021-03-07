package main

// Sample run-helloworld is a minimal Cloud Run service.

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"

	// Import godotenv
	"github.com/joho/godotenv"
)

type WebPage struct {
	Title   string
	Content string
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/slack", slackTestHandler)
	http.HandleFunc("/api/slack/action-endpoint", slackActionEndpoint)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}

	page := WebPage{Title: "Welcome", Content: "Hellow world"}
	t, err := template.ParseFiles("basictemplate.html")

	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, page)
}

func slackTestHandler(w http.ResponseWriter, r *http.Request) {
	slackToken := goDotEnvVariable("SLACK_API_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK API TOKEN MISSING")
	}
	api := slack.New(slackToken)
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}

	channelID, timestamp, err := api.PostMessage(
		"YOUR_CHANNEL_ID",
		slack.MsgOptionText("Some text", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	fmt.Fprintf(w, "slack will talk to this one")
}

func slackActionEndpoint(w http.ResponseWriter, r *http.Request) {
	// Verifies ownership of an Events API Request URL
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		challenge := r.FormValue("challenge")
		w.Write([]byte(challenge))
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
