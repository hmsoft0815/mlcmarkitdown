package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestLlmUtilsHandler_GetTool(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewLlmUtilsHandler(uc)
	tool := h.GetTool()

	if tool.Name != "markitdown__list_models__mlc" {
		t.Errorf("Expected tool name 'markitdown__list_models__mlc', got '%s'", tool.Name)
	}
	if tool.Description == "" {
		t.Error("Expected non-empty description for list_models tool")
	}
}

func TestLlmUtilsHandler_Handle_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			http.NotFound(w, r)
			return
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-api-key" {
			t.Errorf("Expected Authorization 'Bearer test-api-key', got '%s'", auth)
		}

		resp := openAiModelsResponse{
			Data: []ModelInfo{
				{ID: "llama3.2-vision"},
				{ID: "llava:7b"},
				{ID: "mistral"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewLlmUtilsHandler(uc)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markitdown__list_models__mlc",
			Arguments: map[string]interface{}{
				"base_url": mockServer.URL,
				"api_key":  "test-api-key",
			},
		},
	}

	result, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error from Handle: %v", err)
	}
	if result.IsError {
		t.Fatalf("Expected successful result, got tool error")
	}

	if len(result.Content) == 0 {
		t.Fatal("Expected at least one content item")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "llama3.2-vision") {
		t.Errorf("Expected response to contain 'llama3.2-vision', got: %s", textContent.Text)
	}
	if !strings.Contains(textContent.Text, "llava:7b") {
		t.Errorf("Expected response to contain 'llava:7b', got: %s", textContent.Text)
	}
}

func TestLlmUtilsHandler_Handle_ServerError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewLlmUtilsHandler(uc)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markitdown__list_models__mlc",
			Arguments: map[string]interface{}{
				"base_url": mockServer.URL,
			},
		},
	}

	result, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Fatal("Expected tool error result when server returns 500")
	}
}

func TestLlmUtilsHandler_Handle_InvalidJson(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("invalid-json"))
	}))
	defer mockServer.Close()

	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewLlmUtilsHandler(uc)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markitdown__list_models__mlc",
			Arguments: map[string]interface{}{
				"base_url": mockServer.URL,
			},
		},
	}

	result, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected Go error: %v", err)
	}
	if !result.IsError {
		t.Fatal("Expected tool error result when response is invalid JSON")
	}
}
