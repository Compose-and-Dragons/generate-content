package main

import (
	"context"
	"fmt"
	"os"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
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

	systemInstruction, err := helpers.ReadTextFile("instructions.md")
	if err != nil {
		panic(err)
	}

	steps, err := helpers.ReadTextFile("steps.md")
	if err != nil {
		panic(err)
	}

	userContent := "Generate monster data for the given kind: " + kind + " with the above instructions."

	responseFormat := openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONObject: &openai.ResponseFormatJSONObjectParam{
			Type: "json_object",
		},
	}

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       modelRunnerChatModel,
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemInstruction),
				openai.SystemMessage(steps),
				openai.UserMessage(userContent),
			},
			ResponseFormat: responseFormat,

		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelError),
	)
	if err != nil {
		panic(err)
	}
	answer, err := bob.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = helpers.WriteTextFile("contents/monster_sheet.json", answer)
	if err != nil {
		panic(fmt.Errorf("failed to write monster sheet: %w", err))
	}

}
