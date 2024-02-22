package supporthistorystore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

type postgres struct {
	dbConn *pgx.Conn
}

var _ Store = &postgres{}

func NewPostgresStore(dbConn *pgx.Conn) Store {
	return &postgres{dbConn: dbConn}
}

func (store *postgres) FindSimilarSupportCase(ctx context.Context, embeddings []float32) (string, error) {
	row := store.dbConn.QueryRow(
		ctx,
		"SELECT chat FROM support_history ORDER BY embedding <=> $1 LIMIT 1",
		pgvector.NewVector(embeddings),
	)

	var similarSupportCase string
	err := row.Scan(&similarSupportCase)
	if err != nil {

		return "", err
	}

	return similarSupportCase, nil
}

func (store *postgres) SaveSupportCase(ctx context.Context, chat string, embeddings []float32) error {
	_, err := store.dbConn.Exec(
		ctx,
		"INSERT INTO support_history (chat, embedding) VALUES ($1, $2)",
		chat,
		pgvector.NewVector(embeddings),
	)

	return err
}
