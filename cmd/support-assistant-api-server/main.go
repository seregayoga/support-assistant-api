package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
	"github.com/seregayoga/support-assistant-api/internal/handler"
	"github.com/seregayoga/support-assistant-api/internal/llm"
	"github.com/seregayoga/support-assistant-api/internal/supporthistorystore"
)

const (
	dbConnString = "postgres://support-assistant-api:very-secure@localhost:5432/support-assistant-api"
)

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	dbConn, err := pgx.Connect(ctx, dbConnString)
	if err != nil {
		return err
	}
	defer dbConn.Close(ctx)

	supportHistoryStore := supporthistorystore.NewPostgresStore(dbConn)

	openaiToken := os.Getenv("OPENAI_API_KEY")
	if openaiToken == "" {
		return errors.New("Please specify OPENAI_API_KEY env variable")
	}

	openaiClient := openai.NewClient(openaiToken)
	llmClient := llm.NewOpenaiLLM(openaiClient)

	supportHandler := handler.NewSupportHandler(supportHistoryStore, llmClient)

	mux := http.NewServeMux()

	mux.Handle("POST /v1/support", supportHandler)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	log.Println("Listening on:", srv.Addr)

	<-ctx.Done()
	log.Println("interrupted")

	shutdownContext, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return srv.Shutdown(shutdownContext)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
