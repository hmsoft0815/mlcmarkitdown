package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
)

type LlmUtilsHandler struct {
	useCase *usecase.ConvertUseCase
}

func NewLlmUtilsHandler(useCase *usecase.ConvertUseCase) *LlmUtilsHandler {
	return &LlmUtilsHandler{
		useCase: useCase,
	}
}

type ModelInfo struct {
	ID      string `json:"id"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type ModelListResponse struct {
	Models []string `json:"models" jsonschema:"description=List of available model IDs"`
}

type openAiModelsResponse struct {
	Data []ModelInfo `json:"data"`
}

func (h *LlmUtilsHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__list_models__mlc",
		mcp.WithDescription("Lists available models from an OpenAI-compatible LLM provider (like Ollama or OpenAI)."),
		mcp.WithString("base_url", mcp.Description("Optional: The base URL of the LLM API (e.g. http://localhost:11434/v1). If not provided, defaults are used.")),
		mcp.WithString("api_key", mcp.Description("Optional: The API key/token if required.")),
		mcp.WithOutputSchema[ModelListResponse](),
	)
}

func (h *LlmUtilsHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, _, defUrl, defKey := h.useCase.GetLlmDefaults()

	baseUrl := mcp.ParseString(request, "base_url", defUrl)
	if baseUrl == "" {
		// Try to fallback to OpenAI default if nothing else
		baseUrl = "https://api.openai.com/v1"
	}
	
	// Ensure no trailing slash
	baseUrl = strings.TrimSuffix(baseUrl, "/")
	
	apiKey := mcp.ParseString(request, "api_key", defKey)

	req, err := http.NewRequestWithContext(ctx, "GET", baseUrl+"/models", nil)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to connect to LLM provider", err), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("LLM provider returned status %d", resp.StatusCode)), nil
	}

	var modelsResp openAiModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to decode response", err), nil
	}

	modelIds := make([]string, 0, len(modelsResp.Data))
	for _, m := range modelsResp.Data {
		modelIds = append(modelIds, m.ID)
	}

	result := &ModelListResponse{Models: modelIds}
	jsonBytes, _ := json.MarshalIndent(result, "", "  ")

	return mcp.NewToolResultText(fmt.Sprintf("Available Models at %s:\n%s", baseUrl, string(jsonBytes))), nil
}
