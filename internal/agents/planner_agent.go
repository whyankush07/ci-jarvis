package agents

import (
	"context"
	"fmt"
	"log"
	"strings"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/store/vector"
	"ci-jarvis/internal/types"
)

type PlannerAgent struct {
	BaseAgent
	vectorStore vector.VectorStore
}

func NewPlannerAgent(llmClient *llm.Client, vectorStore vector.VectorStore) *PlannerAgent {
	return &PlannerAgent{
		BaseAgent:   NewBaseAgent(llmClient),
		vectorStore: vectorStore,
	}
}

func (p *PlannerAgent) Name() string {
	return "Planner"
}

func (p *PlannerAgent) Execute(ctx context.Context, run *types.Run) (*types.AgentResult, error) {
	log.Printf("[%s] Starting planning for Run ID: %s", p.Name(), run.ID)
	p.fetchAndAttachContext(ctx, run)
	prompt, err := p.LoadPrompt(p.Name(), run)
	if err != nil {
		return nil, fmt.Errorf("failed to load planner prompt: %w", err)
	}

	response, err := p.LLM.GenerateText(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("planner failed to generate text: %w", err)
	}

	var plan types.PlannerOutput
	err = p.ParseJSONResponse(response, &plan)
	if err != nil {
		return nil, fmt.Errorf("failed to parse planner response: %w", err)
	}

	return &types.AgentResult{
		AgentName: p.Name(),
		Success:   true,
		Output:    response,
		Data:      plan,
	}, nil
}

func (p *PlannerAgent) fetchAndAttachContext(ctx context.Context, run *types.Run) {
	if run.Diff != "" && p.vectorStore != nil {
		log.Printf("[%s] Searching for relevant context in vector store...", p.Name())

		queryVector, err := p.LLM.EmbedText(ctx, run.Diff)
		if err != nil {
			log.Printf("warning: failed to embed diff for RAG: %v", err)
		} else {
			// Search vector store for top 3 related snippets
			hits, err := p.vectorStore.Search(ctx, queryVector, 3)
			if err != nil {
				log.Printf("warning: vector search failed: %v", err)
			} else {
				var contextParts []string
				for _, hit := range hits {
					if content, ok := hit["content"].(string); ok {
						path, _ := hit["path"].(string)
						contextParts = append(contextParts, fmt.Sprintf("File: %s\n```go\n%s\n```", path, content))
					}
				}
				if len(contextParts) > 0 {
					// Add retrieved context to the run metadata so it's available to the prompt template
					run.Metadata = map[string]interface{}{
						"RelatedContext": strings.Join(contextParts, "\n\n---\n\n"),
					}
					log.Printf("[%s] Found %d relevant code snippets", p.Name(), len(contextParts))
				}
			}
		}
	}
}
