package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	gpt "github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func main() {
	var choice int
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("no key")
		os.Exit(1)
	}
	fmt.Println("Please choose from the options below (or 0 to quit):\n1) Text Completion\n2) CLI ChatGPT Tool\n3) Identify Libraries In the Code Provided\n4) Generate An Image Using DALL-E")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		ModifiedOpenAIPackage(apiKey)
	case 2:
		CommandLineInterfaceTool(apiKey)
	case 3:
		GetLibrariesFromCode(apiKey)
	case 4:
		ImageCreatorDallE(apiKey)
	case 0:
		fmt.Println("Thank you!")
		os.Exit(1)
	default:
		fmt.Println("Please choose from the options below:\n1) Text Completion\n2)CLI ChatGPT Tool\n3)Identify Libraries In the Code Provided\n4) Generate An Image Using DALL-E")
		fmt.Scanln(&choice)
	}

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
	fmt.Printf("\nPrompt to be completed is: %s", "computers are made of")
	ctx := context.Background()
	client := gpt.NewClient(apiKey)
	resp, err := client.Completion(ctx, gpt.CompletionRequest{
		Prompt:    []string{"computers are made of"},
		MaxTokens: gpt.IntPtr(512),
		Stop:      []string{"."},
		Echo:      true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("\n", resp.Choices[0].Text)
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

func ImageCreatorDallE(apiKey string) {
	c := openai.NewClient(apiKey)
	ctx := context.Background()
	var prompt string
	fmt.Println("Enter a textual description to generate an image from: ")
	fmt.Scanln(&prompt)
	// Sample image by link
	reqUrl := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	}

	respUrl, err := c.CreateImage(ctx, reqUrl)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}
	fmt.Println("Image URL:")
	fmt.Println(respUrl.Data[0].URL)

	// Example image as base64
	reqBase64 := openai.ImageRequest{
		Prompt:         "Michael Jordan dunking a football",
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		fmt.Printf("PNG decode error: %v\n", err)
		return
	}

	file, err := os.Create("example.png")
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		fmt.Printf("PNG encode error: %v\n", err)
		return
	}

	fmt.Println("The image was saved as example.png")
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

type ImagePrompt struct {
	Text string `json:"text"`
}

type ImageResponse struct {
	Image string `json:"image"`
}
