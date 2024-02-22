package supporthistorystore

import "context"

type Store interface {
	FindSimilarSupportCase(ctx context.Context, embeddings []float32) (string, error)
	SaveSupportCase(ctx context.Context, chat string, embeddings []float32) error
}
