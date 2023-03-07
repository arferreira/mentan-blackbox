package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAiResponse struct {
    Choices []struct {
        Text string `json:"text"`
    } `json:"choices"`
}

type OpenAiPrompt struct {
    Prompt      string  `json:"prompt"`
    Temperature float64 `json:"temperature"`
    MaxTokens   int     `json:"max_tokens"`
}

func SendPrompt(prompt string) (string, error) {
    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)
    resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return err.Error(), nil
	}

	fmt.Println(resp.Choices[0].Message.Content)

    return resp.Choices[0].Message.Content, nil
}
