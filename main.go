package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gpt "github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("no key")
		os.Exit(1)
	}
	// OpenAI Package
	OriginalPackageOpenAI(apiKey)
	// Modified OpenAI Package
	ModifiedOpenAIPackage(apiKey)
}
func OriginalPackageOpenAI(apiKey string) {
	c := openai.NewClient(apiKey)
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 10,
		Prompt:    "I am",
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	fmt.Println("Original package output:")
	fmt.Println(resp.Choices[0].Text)
}
func ModifiedOpenAIPackage(apiKey string) {
	fmt.Println("Modified package output: ")
	ctx := context.Background()
	client := gpt.NewClient(apiKey)
	resp, err := client.Completion(ctx, gpt.CompletionRequest{
		Prompt:    []string{"the one thing you need to know about golang programming language is"},
		MaxTokens: gpt.IntPtr(50),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Choices[0].Text)
}
