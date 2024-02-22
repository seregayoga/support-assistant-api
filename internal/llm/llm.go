package llm

import "context"

type LlmClient interface {
	CreateEmbeddings(ctx context.Context, input string) ([]float32, error)
	AnswerUserRequest(ctx context.Context, supportRequest, similarSupportCase string) (string, error)
}
