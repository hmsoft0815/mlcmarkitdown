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
3. You can choose between 'openai' (default) and 'ollama' via the 'llm_provider' argument.
4. For Ollama, you can specify 'ollama_model' (e.g. 'llama3.2-vision') and 'ollama_url' (e.g. 'http://macstudio.fritz.box:11434/v1').
5. For very large documents, the tool will automatically save the result to an artifact storage and return a preview plus an Artifact ID. 
6. If an Artifact ID is returned, inform the user that the full document is available in the artifact storage.
7. If you need to inspect a file's metadata first without converting it, use 'markitdown__quick_inspect__mlc'.

Example workflow:
User: "Describe this image using my local Ollama: cat.jpg"
You: (Call markitdown__convert__mlc with uri="cat.jpg", enable_vision=true, llm_provider="ollama", ollama_url="http://macstudio.fritz.box:11434/v1", ollama_model="llama3.2-vision")
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
