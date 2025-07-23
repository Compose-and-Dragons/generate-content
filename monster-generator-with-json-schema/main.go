package main

import (
	"context"
	"fmt"
	"os"

	"monster-generator/helpers"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
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

	ctx := context.Background()

	clientEngine := openai.NewClient(
		option.WithBaseURL(modelRunnerBaseUrl),
		option.WithAPIKey(""),
	)

	userContent := "Generate monster data for the given kind: " + kind + " with the above instructions."

	responseFormat := openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
			Type: "json_schema",
			JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        "monster_character",
				Description: openai.String("A D&D monster/character data structure"),
				Schema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"Kind": map[string]any{
							"type":        "string",
							"description": "The kind of character to generate",
						},
						"Size": map[string]any{
							"type":        "string",
							"description": "The size of the character to generate",
						},
						"Type": map[string]any{
							"type":        "string",
							"description": "The type of character to generate",
						},
						"Description": map[string]any{
							"type":        "string",
							"description": "The description of the character to generate",
						},
						"Alignment": map[string]any{
							"type":        "string",
							"description": "The alignment of the character to generate",
						},
						"ArmorClass": map[string]any{
							"type":        "integer",
							"description": "The armor class of the character to generate",
						},
						"HitPoints": map[string]any{
							"type":        "integer",
							"description": "The hit points of the character to generate",
						},
						"Speed": map[string]any{
							"type":        "string",
							"description": "The speed of the character to generate",
						},
						"Skills": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "The skills of the character to generate",
						},
						"Senses": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "The senses of the character to generate",
						},
						"Languages": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "The languages of the character to generate",
						},
						"Challenge": map[string]any{
							"type":        "string",
							"description": "The challenge rating of the character to generate",
						},
						"Actions": map[string]any{
							"type":        "array",
							"items":       map[string]any{"type": "string"},
							"description": "The actions of the character to generate",
						},
					},
					"required":             []string{"Kind", "Size", "Type", "Description", "Alignment", "ArmorClass", "HitPoints", "Speed", "Skills", "Senses", "Languages", "Challenge", "Actions"},
					"additionalProperties": false,
				},
			},
		},
	}

	// Chat Completion parameters
	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemInstructions),
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
