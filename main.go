package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	//"testing/quick"

	gpt "github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
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
	// Libraries Detector
	GetLibrariesFromCode(apiKey)
	// CLI Tool
	CommandLineInterfaceTool(apiKey)
}
func OriginalPackageOpenAI(apiKey string) {
	c := openai.NewClient(apiKey)
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 10,
		Prompt:    "I am",
		Echo:      true,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	fmt.Println("Original package output:")
	fmt.Print(resp.Choices[0].Text)
}
func ModifiedOpenAIPackage(apiKey string) {
	fmt.Println("\nModified package output: ")
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
	fmt.Print(resp.Choices[0].Text)
}
func CommandLineInterfaceTool(apiKey string) {
	fmt.Println("\nCLI Tool:")
	log.SetOutput(new(NullWriter))
	ctx := context.Background()
	client := gpt.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Println("Say something ('quit' to end):")
				if !scanner.Scan() {
					break
				}
				question := scanner.Text()
				switch question {
				case "quit":
					quit = true
				default:
					GetResponse(ctx, client, question)
				}
			}
		},
	}
	rootCmd.Execute()
}
func GetResponse(ctx context.Context, client gpt.Client, question string) {
	err := client.CompletionStreamWithEngine(ctx, gpt.TextDavinci003Engine, gpt.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt.IntPtr(512),
		Temperature: gpt.Float32Ptr(0),
	}, func(resp *gpt.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
	fmt.Printf("\n")

}
func GetLibrariesFromCode(apiKey string) {
	fmt.Println("\nLibraries Detector Tool:")
	log.SetOutput(new(NullWriter))
	ctx := context.Background()
	client := gpt.NewClient(apiKey)

	const inputFile = "./input_with_code.txt"
	fileBytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	// gpt prompt - backticks to indicate a code snippet for gpt interpreter
	msgPrefix := "give me a shortlist of the libraries that are used in the code \n``` python\n"
	msgSuffix := "\n```"
	msg := msgPrefix + string(fileBytes) + msgSuffix

	outputBuilder := strings.Builder{}

	err = client.CompletionStreamWithEngine(ctx, gpt.TextDavinci003Engine, gpt.CompletionRequest{
		Prompt: []string{
			msg,
		},
		MaxTokens:   gpt.IntPtr(3000),
		Temperature: gpt.Float32Ptr(0),
	}, func(resp *gpt.CompletionResponse) {
		outputBuilder.WriteString(resp.Choices[0].Text)
	})
	if err != nil {
		log.Fatalln(err)
	}
	output := strings.TrimSpace(outputBuilder.String())
	const outputFile = "./output.txt"
	err = os.WriteFile(outputFile, []byte(output), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }
