package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/types"
)

type PlannerAgent struct {
	BaseAgent
}

func NewPlannerAgent(llmClient *llm.Client) *PlannerAgent {
	return &PlannerAgent{
		BaseAgent: NewBaseAgent(llmClient),
	}
}

func (p *PlannerAgent) Name() string {
	return "Planner"
}

// 1. Fetch PR metadata and diff.
// 2. Send a prompt to the LLM with this context.
// 3. Parse the LLM's suggested steps into a structured Plan.

func (p *PlannerAgent) Execute(ctx context.Context, run *types.Run) (*types.AgentResult, error) {
	log.Printf("[%s] Starting planning for Run ID: %s", p.Name(), run.ID)

	// Load and execute the prompt template
	prompt, err := p.LoadPrompt(p.Name(), run)
	if err != nil {
		return nil, fmt.Errorf("failed to load planner prompt: %w", err)
	}

	response, err := p.LLM.GenerateText(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("planner failed to generate text: %w", err)
	}

	// Parse the LLM response into a structured PlannerOutput
	var plan types.PlannerOutput
	err = p.ParseLLMResponse(response, &plan)
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

func (p *PlannerAgent) ParseLLMResponse(response string, plan *types.PlannerOutput) error {
	// Clean up LLM response (strip markdown code blocks if present)
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

	return json.Unmarshal([]byte(cleanJSON), plan)
}
