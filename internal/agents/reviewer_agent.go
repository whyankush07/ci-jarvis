package agents

import (
	"context"
	"fmt"
	"log"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/types"
)

type ReviewerAgent struct {
	BaseAgent
}

func NewReviewerAgent(llmClient *llm.Client) *ReviewerAgent {
	return &ReviewerAgent{
		BaseAgent: NewBaseAgent(llmClient),
	}
}

func (r *ReviewerAgent) Name() string {
	return "Reviewer"
}

func (r *ReviewerAgent) Execute(ctx context.Context, run *types.Run) (*types.AgentResult, error) {
	log.Printf("[%s] Starting review for Run ID: %s", r.Name(), run.ID)

	// Load and execute the prompt template
	prompt, err := r.LoadPrompt(r.Name(), run)
	if err != nil {
		return nil, fmt.Errorf("failed to load reviewer prompt: %w", err)
	}

	response, err := r.LLM.GenerateText(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("reviewer failed to generate text: %w", err)
	}

	// Parse the LLM response into a structured ReviewerOutput
	var output types.ReviewerOutput
	err = r.ParseJSONResponse(response, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reviewer response: %w", err)
	}

	return &types.AgentResult{
		AgentName: r.Name(),
		Success:   true,
		Output:    response,
		Data:      output,
	}, nil
}
