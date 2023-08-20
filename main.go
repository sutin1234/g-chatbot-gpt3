package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"github.com/sutin1234/go-chatbot-gpt3/chatbot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("CHAT_GPT_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing API_KEY in .env file")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	resp, err := client.Completion(ctx, gpt3.CompletionRequest{
		Prompt:    []string{"you should know about golang is"},
		MaxTokens: gpt3.IntPtr(30),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices[0].Text)

	// check env
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		log.Fatal("supabaseURL missing in .env")
	}

	supabaseAPIKey := os.Getenv("SUPABASE_API_KEY")
	if supabaseAPIKey == "" {
		log.Fatal("supabaseAPIKey missing in .env")
	}

	chatGptCompletionURL := os.Getenv("CHAT_GPT_COMPLETION_URL")
	if chatGptCompletionURL == "" {
		log.Fatal("chatGptCompletionURL missing in .env")
	}

	text, err := chatbot.AutoCompletions("Hello Chat GPT")
	if err != nil {
		log.Fatal("Cannot fetch autocompletion", err)
	}
	fmt.Printf("Chatbot: %s", text)
}
