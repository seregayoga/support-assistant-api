package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/seregayoga/support-assistant-api/internal/llm"
	"github.com/seregayoga/support-assistant-api/internal/supporthistorystore"
)

type supportRequest struct {
	SupportRequest string `json:"support_request"`
}

type SupportHandler struct {
	supportHistoryStore supporthistorystore.Store
	llmClient           llm.LlmClient
}

func NewSupportHandler(supportHistoryStore supporthistorystore.Store, llmClient llm.LlmClient) http.Handler {
	return &SupportHandler{
		supportHistoryStore: supportHistoryStore,
		llmClient:           llmClient,
	}
}

func (handler *SupportHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	var request supportRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Printf("Failed to parse request: %s", err)
		http.Error(w, "Failed to parse request", http.StatusBadRequest)

		return
	}

	embeddings, err := handler.llmClient.CreateEmbeddings(req.Context(), request.SupportRequest)
	if err != nil {
		log.Printf("Failed to get embeddings: %s", err)
		http.Error(w, "Failed to get embeddings", http.StatusInternalServerError)

		return
	}

	similarSupportCase, err := handler.supportHistoryStore.FindSimilarSupportCase(req.Context(), embeddings)
	if err != nil {
		log.Printf("Failed to get chat from db: %s", err)
		http.Error(w, "Failed to get chat from db", http.StatusInternalServerError)

		return
	}

	answer, err := handler.llmClient.AnswerUserRequest(req.Context(), request.SupportRequest, similarSupportCase)
	if err != nil {
		log.Printf("Failed to get answer: %s", err)
		http.Error(w, "Failed to get answer", http.StatusInternalServerError)

		return
	}

	response := struct {
		Answer string `json:"answer"`
	}{Answer: answer}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed write response: %s", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
