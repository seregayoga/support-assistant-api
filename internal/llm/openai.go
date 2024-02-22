package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type openaiLLM struct {
	openaiClient *openai.Client
}

var _ LlmClient = &openaiLLM{}

func NewOpenaiLLM(openaiClient *openai.Client) LlmClient {
	return &openaiLLM{
		openaiClient: openaiClient,
	}
}

func (openaiLLM *openaiLLM) CreateEmbeddings(ctx context.Context, input string) ([]float32, error) {
	request := openai.EmbeddingRequest{
		Input: input,
		Model: openai.SmallEmbedding3,
	}

	response, err := openaiLLM.openaiClient.CreateEmbeddings(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Data[0].Embedding, err
}

func (llm *openaiLLM) AnswerUserRequest(ctx context.Context, supportRequest, similarSupportCase string) (string, error) {
	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant providing support to users of a local internet provider.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: supportRequest,
			},
			{
				Role: openai.ChatMessageRoleAssistant,
				Content: fmt.Sprintf(
					`Context: 
The user's current request is similar to a previously solved support case. Here is the information from the relevant historical case:

Historical Support Case:
"%s"

Provide direct instructions to the user based on the insights gained from the historical support case.`,
					similarSupportCase,
				),
			},
		},
	}

	response, err := llm.openaiClient.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
