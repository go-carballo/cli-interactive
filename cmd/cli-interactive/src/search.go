package agent

import (
	"context"
	"fmt"

	tavily "github.com/hekmon/tavily/v2"
)

// SearchWeb searches the web using Tavily API
func SearchWeb(ctx context.Context, client tavily.Client, query string, maxResults int) (tavily.SearchAnswer, error) {
	if maxResults <= 0 {
		maxResults = 5
	}

	response, err := client.Search(ctx, tavily.SearchQuery{
		Query:             query,
		SearchDepth:       tavily.SearchQueryDepthAdvanced,
		MaxResults:        maxResults,
		IncludeAnswer:     tavily.SearchQueryIncludeAnswerBasic,
		IncludeRawContent: false,
		IncludeImages:     false,
	})
	if err != nil {
		return tavily.SearchAnswer{}, fmt.Errorf("Tavily search failed: %w", err)
	}

	return response, nil
}
