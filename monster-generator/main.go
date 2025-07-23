package main

import (
	"context"
	"fmt"
	"os"

	"monster-generator/helpers"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")

	if modelRunnerBaseUrl == "" {
		panic("MODEL_RUNNER_BASE_URL environment variable is not set")
	}
	modelRunnerChatModel := os.Getenv("MODEL_RUNNER_CHAT_MODEL")
	fmt.Println("Using Model Runner Chat Model:", modelRunnerChatModel)

	if modelRunnerChatModel == "" {
		panic("MODEL_RUNNER_CHAT_MODEL environment variable is not set")
	}

	kind := os.Getenv("MONSTER_KIND")
	if kind == "" {
		panic("MONSTER_KIND environment variable is not set")
	}

	systemInstructions, err := helpers.ReadTextFile("instructions.md")
	if err != nil {
		panic(err)
	}

	steps, err := helpers.ReadTextFile("steps.md")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	clientEngine := openai.NewClient(
		option.WithBaseURL(modelRunnerBaseUrl),
		option.WithAPIKey(""),
	)

	userContent := "Generate monster data for the given kind: " + kind + " with the above instructions."

	responseFormat := openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONObject: &openai.ResponseFormatJSONObjectParam{
			Type: "json_object",
		},
	}

	// Chat Completion parameters
	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemInstructions),
			openai.SystemMessage(steps),
			openai.UserMessage(userContent),
		},
		ResponseFormat: responseFormat,
		Model:          modelRunnerChatModel,
		Temperature:    openai.Opt(0.8),
	}

	stream := clientEngine.Chat.Completions.NewStreaming(ctx, chatCompletionParams)

	answer := ""
	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			answer += chunk.Choices[0].Delta.Content
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("😡 Stream error: %v\n", err)
	}

	err = helpers.WriteTextFile("contents/monster_sheet.json", answer)
	if err != nil {
		panic(fmt.Errorf("failed to write monster sheet: %w", err))
	}

}
