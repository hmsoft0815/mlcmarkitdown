package handler

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type PromptHandler struct{}

func NewPromptHandler() *PromptHandler {
	return &PromptHandler{}
}

func (h *PromptHandler) GetPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"markitdown_instruction",
		mcp.WithPromptDescription("Instructions for the LLM on how to use the mlc-markitdown tools for document conversion."),
	)
}

func (h *PromptHandler) Handle(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	instruction := `You are an expert at document conversion using the mlc-markitdown toolset. 

When a user asks to see the content of a file (PDF, DOCX, XLSX, PPTX, HTML, etc.) or a URL in Markdown format, you should use the tool 'markitdown__convert__mlc'.

Key Instructions:
1. Always provide the 'uri' argument (local path or URL).
2. For images (JPG, PNG) or audio (MP3, WAV), set 'enable_vision' to true.
3. The server is pre-configured with a default LLM provider (e.g. Ollama or OpenAI). You usually don't need to specify 'llm_provider', 'ollama_url', etc., unless you want to override the defaults.
4. If you are unsure which models are available, use the 'markitdown__list_models__mlc' tool.
5. For very large documents, the tool will automatically save the result to an artifact storage and return a preview plus an Artifact ID. 
6. If an Artifact ID is returned, inform the user that the full document is available in the artifact storage.
7. If you need to inspect a file's metadata first without converting it, use 'markitdown__quick_inspect__mlc'.

Example workflow:
User: "Describe this image: cat.jpg"
You: (Call markitdown__convert__mlc with uri="cat.jpg", enable_vision=true)
`
	return mcp.NewGetPromptResult(
		"MLC MarkItDown Usage Instructions",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleAssistant,
				mcp.NewTextContent(instruction),
			),
		},
	), nil
}
