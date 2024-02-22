package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
	"github.com/seregayoga/support-assistant-api/internal/llm"
	"github.com/seregayoga/support-assistant-api/internal/supporthistorystore"
)

const (
	supportHistoryLogNewFolder      = "./support_history_log/new/"
	supportHistoryLogMigratedFolder = "./support_history_log/migrated/"
	dbConnString                    = "postgres://support-assistant-api:very-secure@localhost:5432/support-assistant-api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbConn, err := pgx.Connect(ctx, dbConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(ctx)
	supportHistoryStore := supporthistorystore.NewPostgresStore(dbConn)

	openaiToken := os.Getenv("OPENAI_API_KEY")
	openaiClient := openai.NewClient(openaiToken)
	llmClient := llm.NewOpenaiLLM(openaiClient)

	files, err := os.ReadDir(supportHistoryLogNewFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println("Migrating", file.Name())

		filePath := supportHistoryLogNewFolder + file.Name()

		chatBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		chat := string(chatBytes)

		embeddings, err := llmClient.CreateEmbeddings(ctx, chat)
		if err != nil {
			log.Fatal(err)
		}

		err = supportHistoryStore.SaveSupportCase(ctx, chat, embeddings)
		if err != nil {
			log.Fatal(err)
		}

		migratedFilePath := supportHistoryLogMigratedFolder + file.Name()
		err = os.Rename(filePath, migratedFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
