# CLI-Interactive

> An intelligent command-line AI assistant powered by Google Gemini and real-time web search capabilities.

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Gemini](https://img.shields.io/badge/Gemini_2.5_Pro-4285F4?style=flat&logo=google&logoColor=white)](https://ai.google.dev/)
[![Firebase Genkit](https://img.shields.io/badge/Firebase_Genkit-FFCA28?style=flat&logo=firebase&logoColor=black)](https://firebase.google.com/docs/genkit)
[![Tavily](https://img.shields.io/badge/Tavily_API-Search-green?style=flat)](https://tavily.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Overview

**CLI-Interactive** is a conversational AI agent that runs directly in your terminal. It combines the power of **Google Gemini 2.5 Pro** with **real-time web search** through Tavily API, allowing you to get accurate, up-to-date answers with cited sources.

### Key Features

- **Conversational AI** - Natural language interaction with chat history maintained across the session
- **Real-time Web Search** - Integrated Tavily API for fetching current information from the web
- **Source Citations** - All answers include numbered references with URLs for verification
- **Interactive CLI** - Beautiful terminal interface with colors, spinners, and readline support
- **Genkit Integration** - Built on Firebase Genkit for robust AI orchestration and tool calling
- **Dual Mode** - Run as interactive CLI or as a Genkit server for UI integration

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI-Interactive                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      │
│  │   Readline  │───▶│   Agent     │───▶│   Gemini    │      │
│  │  Interface  │    │  (Genkit)   │    │   2.5 Pro   │      │
│  └─────────────┘    └──────┬──────┘    └─────────────┘      │
│                            │                                 │
│                     ┌──────▼──────┐                         │
│                     │  searchWeb  │                         │
│                     │    Tool     │                         │
│                     └──────┬──────┘                         │
│                            │                                 │
│                     ┌──────▼──────┐                         │
│                     │ Tavily API  │                         │
│                     └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

---

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Google Gemini API key ([Get one here](https://aistudio.google.com/apikey))
- Tavily API key ([Get one here](https://tavily.com/))

### Installation

```bash
# Clone the repository
git clone https://github.com/go-carballo/cli-interactive.git
cd cli-interactive

# Install dependencies
go mod download

# Build the binary
go build -o cli-interactive ./cmd/cli-interactive/
```

### Configuration

Create a `.env` file in the project root:

```env
GEMINI_API_KEY=your_gemini_api_key_here
TAVILY_API_KEY=your_tavily_api_key_here
```

### Usage

```bash
# Run in interactive mode
./cli-interactive

# Run in Genkit server mode (for UI integration)
./cli-interactive serve
```

---

## Demo

```
╔════════════════════════════════════════════════════════════╗
║         Welcome to CLI-Interactive - Interactive Mode       ║
╚════════════════════════════════════════════════════════════╝
Type your questions and get AI-powered answers with sources!
Chat history is maintained during this session.

😊 Ask a question (or type "exit" to quit): What are the latest developments in AI?

🤖 Thinking...

💬 Answer:
Based on recent developments, here are the key highlights in AI...

**Sources**
[1] https://example.com/ai-news-1
[2] https://example.com/ai-news-2
```

---

## Commands

| Command | Description |
|---------|-------------|
| `exit` / `quit` | Exit the application |
| `clear` | Clear chat history |
| `Ctrl+C` | Force exit |

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.24** | Core programming language |
| **Firebase Genkit** | AI orchestration framework with tool calling |
| **Google Gemini 2.5 Pro** | Large language model for generation |
| **Tavily API** | Real-time web search engine for AI |
| **Readline** | Terminal input with history and editing |
| **Fatih Color** | Terminal colors and styling |
| **Spinner** | Loading animations |

---

## Project Structure

```
cli-interactive/
├── cmd/
│   └── cli-interactive/
│       ├── main.go           # Application entry point
│       └── src/
│           ├── agent.go      # Chat agent implementation
│           └── search.go     # Tavily search wrapper
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── .env                      # Environment variables (not tracked)
├── .gitignore
└── README.md
```

---

## How It Works

1. **User Input** - The user types a question in the terminal
2. **Agent Processing** - The Genkit agent receives the query with conversation history
3. **Tool Calling** - Gemini decides to call the `searchWeb` tool when current information is needed
4. **Web Search** - Tavily API performs an advanced search and returns relevant results
5. **Response Generation** - Gemini synthesizes the search results into a comprehensive answer
6. **Source Citation** - The response includes numbered citations linking to original sources

---

## Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Author

**Alejandro Carballo**

- GitHub: [@go-carballo](https://github.com/go-carballo)

---

## Tags

`#AI` `#ArtificialIntelligence` `#GenerativeAI` `#LLM` `#Gemini` `#GoogleAI` `#CLI` `#CommandLine` `#Terminal` `#Go` `#Golang` `#Firebase` `#Genkit` `#WebSearch` `#Tavily` `#ChatBot` `#AIAgent` `#ToolCalling` `#FunctionCalling` `#RAG` `#OpenSource` `#Developer` `#Programming` `#Tech` `#Innovation`

---

<p align="center">
  Built with Go and powered by Google Gemini
</p>
