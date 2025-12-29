package agent

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/hekmon/tavily/v2"
)

// SearchInput defines the input schema for the search tool
type SearchInput struct {
	Query string `json:"query" jsonschema:"description=The search query to look up"`
}

// ChatInput defines the input schema for the chat prompt
type ChatInput struct {
	Query string `json:"query" jsonschema:"description=The user query to be answered using web search results"`
}

// Agent represents a chat agent with web search capabilities
type Agent struct {
	g            *genkit.Genkit
	model        ai.Model
	searchTool   ai.Tool
	history      []*ai.Message
	systemPrompt string
}

// CreateSearchTool creates a search tool definition for the chat agent
func CreateSearchTool(g *genkit.Genkit, tavilyClient tavily.Client) ai.Tool {
	return genkit.DefineTool(g, "searchWeb",
		"Search the web for current information to answer user queries. Use this when you need up-to-date or factual information.",
		func(ctx *ai.ToolContext, input SearchInput) (string, error) {
			// Search using Tavily
			results, err := tavilyClient.Search(ctx.Context, tavily.SearchQuery{
				Query:      input.Query,
				MaxResults: 5,
			})
			if err != nil {
				return "", fmt.Errorf("search error: %w", err)
			}

			// Format search results for the model
			var formattedResults string
			for i, result := range results.Results {
				formattedResults += fmt.Sprintf("[%d] %s\nURL: %s\nContent: %s\n\n",
					i+1, result.Title, result.URL, result.Content)
			}

			return formattedResults, nil
		},
	)
}

// CreateChatAgent creates a chat agent with web search capabilities
func CreateChatAgent(g *genkit.Genkit, tavilyClient tavily.Client, model ai.Model) *Agent {
	// Define the search tool
	searchTool := CreateSearchTool(g, tavilyClient)

	systemPrompt := `You are a helpful AI assistant that provides comprehensive and accurate answers based on web search results.

Instructions:
1. Provide a comprehensive answer to the user's query based on the search results above
2. Synthesize information from multiple sources when relevant
3. Be factual and cite specific sources using [1], [2], etc. notation
4. If the search results don't contain enough information, acknowledge this
5. Keep the answer clear and well-structured
6. Use markdown formatting for better readability
7. Please use the tool searchWeb always when you need to look up current information
8. Add a section at the end titled "Sources" listing the URLs of the references used`

	return &Agent{
		g:            g,
		model:        model,
		searchTool:   searchTool,
		history:      []*ai.Message{},
		systemPrompt: systemPrompt,
	}
}

// Send sends a message to the agent and returns the response
func (a *Agent) Send(ctx context.Context, query string) (string, error) {
	// Add user message to history
	a.history = append(a.history, &ai.Message{
		Role:    ai.RoleUser,
		Content: []*ai.Part{ai.NewTextPart(query)},
	})

	// Generate response with tools
	response, err := genkit.Generate(ctx, a.g,
		ai.WithModel(a.model),
		ai.WithSystem(a.systemPrompt),
		ai.WithMessages(a.history...),
		ai.WithTools(a.searchTool),
	)
	if err != nil {
		return "", fmt.Errorf("generate error: %w", err)
	}

	// Get response text
	responseText := response.Text()

	// Add assistant message to history
	a.history = append(a.history, &ai.Message{
		Role:    ai.RoleModel,
		Content: []*ai.Part{ai.NewTextPart(responseText)},
	})

	return responseText, nil
}

// ClearHistory clears the chat history
func (a *Agent) ClearHistory() {
	a.history = []*ai.Message{}
}
