package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/hekmon/tavily/v2"
	"github.com/joho/godotenv"
)

// SearchInput defines the input schema for the search tool
type SearchInput struct {
	Query string `json:"query" jsonschema:"description=The search query to look up"`
}

// Agent represents a chat agent with web search capabilities
type Agent struct {
	g            *genkit.Genkit
	model        ai.Model
	searchTool   ai.Tool
	history      []*ai.Message
	systemPrompt string
}

// createSearchTool creates a search tool definition for the chat agent
func createSearchTool(g *genkit.Genkit, tavilyClient tavily.Client) ai.Tool {
	return genkit.DefineTool(g, "searchWeb",
		"Search the web for current information to answer user queries. Use this when you need up-to-date or factual information.",
		func(ctx *ai.ToolContext, input SearchInput) (string, error) {
			// Search using Tavily
			results, err := tavilyClient.Search(ctx.Context, tavily.SearchQuery{
				Query:             input.Query,
				SearchDepth:       tavily.SearchQueryDepthAdvanced,
				MaxResults:        5,
				IncludeAnswer:     tavily.SearchQueryIncludeAnswerBasic,
				IncludeRawContent: false,
				IncludeImages:     false,
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

// createChatAgent creates a chat agent with web search capabilities
func createChatAgent(g *genkit.Genkit, searchTool ai.Tool, model ai.Model) *Agent {
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

func printWelcome() {
	cyan := color.New(color.FgCyan, color.Bold)
	gray := color.New(color.FgHiBlack)

	cyan.Println("\n╔════════════════════════════════════════════════════════════╗")
	cyan.Println("║         Welcome to CLI-Interactive - Interactive Mode       ║")
	cyan.Println("╚════════════════════════════════════════════════════════════╝")
	gray.Println("Type your questions and get AI-powered answers with sources!")
	gray.Println("Chat history is maintained during this session.")
}

func startInteractive(ctx context.Context, g *genkit.Genkit, searchTool ai.Tool, model ai.Model) error {
	// Crear el agente con búsqueda web
	chatAgent := createChatAgent(g, searchTool, model)

	// Mostrar bienvenida
	printWelcome()

	// Crear interfaz readline
	green := color.New(color.FgGreen, color.Bold)
	rl, err := readline.New(green.Sprint("😊 Ask a question (or type \"exit\" to quit): "))
	if err != nil {
		return fmt.Errorf("readline error: %w", err)
	}
	defer rl.Close()

	// Loop principal
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		input := strings.TrimSpace(line)

		// Comandos de salida
		if input == "exit" || input == "quit" {
			color.Yellow("\n👋 Exiting CLI-Interactive. Goodbye!\n")
			break
		}

		if input == "" {
			gray := color.New(color.FgHiBlack)
			gray.Println("Type a question or use: exit, quit, Ctrl+C")
			continue
		}

		// Comando para limpiar historial
		if input == "clear" {
			chatAgent.ClearHistory()
			color.Green("Chat history cleared!")
			continue
		}

		// Enviar mensaje al agente
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " 🤖 Thinking..."
		s.Start()
		response, err := chatAgent.Send(ctx, input)
		s.Stop()
		if err != nil {
			color.Red("Error: %v", err)
			continue
		}

		color.Cyan("\n💬 Answer:")
		fmt.Println(response)
		fmt.Println()
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		color.Red("\n✗ Error loading .env file")
		os.Exit(1)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	tavilyKey := os.Getenv("TAVILY_API_KEY")

	if geminiKey == "" {
		color.Red("\n✗ Error: GEMINI_API_KEY is not set")
		color.Yellow("\n💡 Tip: Make sure to set your GEMINI_API_KEY in the .env file")
		os.Exit(1)
	}

	ctx := context.Background()

	// Inicializar Genkit con el plugin de Google AI
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: geminiKey}),
	)

	// Crear cliente de Tavily (si está disponible)
	var searchTool ai.Tool
	if tavilyKey != "" {
		tavilyClient, err := tavily.NewClient(tavilyKey, nil)
		if err == nil {
			searchTool = createSearchTool(g, tavilyClient)
		}
	}

	// Definir el modelo una sola vez
	model := googlegenai.GoogleAIModel(g, "gemini-2.5-pro")

	// Definir el prompt con la tool de búsqueda
	type QueryInput struct {
		Query string `json:"query" jsonschema:"description=The user query to search for,minLength=1,maxLength=2000"`
	}

	systemMessage := `You are a helpful AI assistant that provides comprehensive and accurate answers based on web search results.

Instructions:
1. Use the searchWeb tool to find current information when needed
2. Provide a comprehensive answer to the user's query based on the search results
3. Synthesize information from multiple sources when relevant
4. Be factual and cite specific sources using [1], [2], etc. notation
5. If the search results don't contain enough information, acknowledge this
6. Keep the answer clear and well-structured
7. Use markdown formatting for better readability
8. Add a section at the end titled "Sources" listing the URLs of the references used`

	searchPrompt := genkit.DefinePrompt(g, "searchPrompt",
		ai.WithDescription("Prompt that searches the web to answer user queries based on current information."),
		ai.WithInputType(QueryInput{}),
		ai.WithModel(model),
		ai.WithSystem(systemMessage),
		ai.WithPrompt("User Query: {{Query}}"),
		ai.WithTools(searchTool),
	)

	// Definir el flow askQuestion para la UI de Genkit
	genkit.DefineFlow(g, "askQuestion", func(ctx context.Context, input QueryInput) (string, error) {
		response, err := searchPrompt.Execute(ctx,
			ai.WithInput(input),
		)
		if err != nil {
			return "", err
		}
		return response.Text(), nil
	})

	// Verificar si se quiere iniciar en modo Genkit UI
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		color.Cyan("\n🚀 Genkit flows registered. Waiting for genkit CLI...")
		select {}
	}

	// Por defecto, iniciar en modo interactivo CLI
	if tavilyKey == "" || searchTool == nil {
		color.Red("\n✗ Error: TAVILY_API_KEY is not set")
		color.Yellow("\n💡 Tip: Make sure to set your TAVILY_API_KEY in the .env file")
		os.Exit(1)
	}

	if err := startInteractive(ctx, g, searchTool, model); err != nil {
		color.Red("\n✗ Error: %s", err.Error())
		os.Exit(1)
	}
}
