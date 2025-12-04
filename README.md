# Seville Weather MCP Server ☀️

MCP server that checks if it's sunny in Seville, Spain. Works with any MCP client (Claude Desktop, custom clients, etc).

## Setup

1. **Get an API key** from [OpenWeatherMap](https://openweathermap.org/) (free tier)

2. **Build the server:**

```bash
go build -o seville-weather main.go
```

3. **Set environment variable:**

```bash
export OPENWEATHER_API_KEY="your_api_key_here"
```

## Usage

### With Claude Desktop

Edit `claude_desktop_config.json`:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "seville-weather": {
      "command": "/absolute/path/to/seville-weather",
      "env": {
        "OPENWEATHER_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

Restart Claude Desktop and ask: "Is it sunny in Seville?"

### With Custom MCP Clients

Run the server with stdio transport:

```bash
./seville-weather
```

The server exposes one tool: `is_sunny` (no parameters required)

## Tool

**`is_sunny`** - Returns current weather status in Seville with temperature
