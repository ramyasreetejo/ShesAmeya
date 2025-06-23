package model

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	model *genai.GenerativeModel
	ctx   context.Context
}

func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	return &GeminiClient{
		model: client.GenerativeModel("gemini-1.5-flash"),
		ctx:   ctx,
	}
}

func (g *GeminiClient) Generate(prompt string) (string, error) {
	resp, err := g.model.GenerateContent(g.ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	for _, c := range resp.Candidates {
		if len(c.Content.Parts) > 0 {
			if t, ok := c.Content.Parts[0].(genai.Text); ok {
				return string(t), nil
			}
		}
	}
	return "", errors.New("no response generated")
}
