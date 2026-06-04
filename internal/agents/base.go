package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/types"
)

// Agent is the interface that all Jarvis agents must implement.
type Agent interface {
	Name() string
	Execute(ctx context.Context, run *types.Run) (*types.AgentResult, error)
}

// BaseAgent provides common functionality for all agents.
type BaseAgent struct {
	LLM *llm.Client
}

func NewBaseAgent(llmClient *llm.Client) BaseAgent {
	return BaseAgent{
		LLM: llmClient,
	}
}

// LoadPrompt reads a prompt template from disk and executes it with the provided data.
func (b *BaseAgent) LoadPrompt(agentName string, data interface{}) (string, error) {
	// In a real project, this path might be configurable via environment variables
	fileName := strings.ToLower(agentName) + ".txt"
	path := filepath.Join("internal", "llm", "prompts", fileName)

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}

	tmpl, err := template.New(agentName).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template: %w", err)
	}

	return buf.String(), nil
}

// ParseJSONResponse cleans up Markdown-wrapped JSON and unmarshals it into the target struct.
func (b *BaseAgent) ParseJSONResponse(response string, target interface{}) error {
	cleanJSON := response
	if strings.Contains(cleanJSON, "```json") {
		parts := strings.Split(cleanJSON, "```json")
		if len(parts) > 1 {
			cleanJSON = strings.Split(parts[1], "```")[0]
		}
	} else if strings.Contains(cleanJSON, "```") {
		parts := strings.Split(cleanJSON, "```")
		if len(parts) > 1 {
			cleanJSON = strings.Split(parts[1], "```")[0]
		}
	}

	cleanJSON = strings.TrimSpace(cleanJSON)

	return json.Unmarshal([]byte(cleanJSON), target)
}
