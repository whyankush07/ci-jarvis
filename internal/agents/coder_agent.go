package agents

import (
	"context"
	"fmt"
	"log"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/types"
)

type CoderAgent struct {
	BaseAgent
}

func NewCoderAgent(llmClient *llm.Client) *CoderAgent {
	return &CoderAgent{
		BaseAgent: NewBaseAgent(llmClient),
	}
}

func (c *CoderAgent) Name() string {
	return "Coder"
}

func (c *CoderAgent) Execute(ctx context.Context, run *types.Run) (*types.AgentResult, error) {
	log.Printf("[%s] Starting coding for Run ID: %s", c.Name(), run.ID)

	// Load and execute the prompt template
	prompt, err := c.LoadPrompt(c.Name(), run)
	if err != nil {
		return nil, fmt.Errorf("failed to load coder prompt: %w", err)
	}

	response, err := c.LLM.GenerateText(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("coder failed to generate text: %w", err)
	}

	// Parse the LLM response into a structured CoderOutput
	var output types.CoderOutput
	err = c.ParseJSONResponse(response, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coder response: %w", err)
	}

	return &types.AgentResult{
		AgentName: c.Name(),
		Success:   true,
		Output:    response,
		Data:      output,
	}, nil
}
